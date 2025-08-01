package scanner

import (
	"fmt"
	"iter"
	"maps"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/core"
	"github.com/microsoft/typescript-go/diagnostics"
	"github.com/microsoft/typescript-go/jsnum"
	"github.com/microsoft/typescript-go/stringutil"
)

type EscapeSequenceScanningFlags int32

const (
	EscapeSequenceScanningFlagsString                     EscapeSequenceScanningFlags = 1 << 0
	EscapeSequenceScanningFlagsReportErrors               EscapeSequenceScanningFlags = 1 << 1
	EscapeSequenceScanningFlagsRegularExpression          EscapeSequenceScanningFlags = 1 << 2
	EscapeSequenceScanningFlagsAnnexB                     EscapeSequenceScanningFlags = 1 << 3
	EscapeSequenceScanningFlagsAnyUnicodeMode             EscapeSequenceScanningFlags = 1 << 4
	EscapeSequenceScanningFlagsAtomEscape                 EscapeSequenceScanningFlags = 1 << 5
	EscapeSequenceScanningFlagsReportInvalidEscapeErrors  EscapeSequenceScanningFlags = EscapeSequenceScanningFlagsRegularExpression | EscapeSequenceScanningFlagsReportErrors
	EscapeSequenceScanningFlagsAllowExtendedUnicodeEscape EscapeSequenceScanningFlags = EscapeSequenceScanningFlagsString | EscapeSequenceScanningFlagsAnyUnicodeMode
)

type ErrorCallback func(diagnostic *diagnostics.Message, start, length int, args ...any)

var textToKeyword = map[string]ast.Kind{
	"abstract":    ast.KindAbstractKeyword,
	"accessor":    ast.KindAccessorKeyword,
	"any":         ast.KindAnyKeyword,
	"as":          ast.KindAsKeyword,
	"asserts":     ast.KindAssertsKeyword,
	"assert":      ast.KindAssertKeyword,
	"bigint":      ast.KindBigIntKeyword,
	"boolean":     ast.KindBooleanKeyword,
	"break":       ast.KindBreakKeyword,
	"case":        ast.KindCaseKeyword,
	"catch":       ast.KindCatchKeyword,
	"class":       ast.KindClassKeyword,
	"continue":    ast.KindContinueKeyword,
	"const":       ast.KindConstKeyword,
	"constructor": ast.KindConstructorKeyword,
	"debugger":    ast.KindDebuggerKeyword,
	"declare":     ast.KindDeclareKeyword,
	"default":     ast.KindDefaultKeyword,
	"delete":      ast.KindDeleteKeyword,
	"do":          ast.KindDoKeyword,
	"else":        ast.KindElseKeyword,
	"enum":        ast.KindEnumKeyword,
	"export":      ast.KindExportKeyword,
	"extends":     ast.KindExtendsKeyword,
	"false":       ast.KindFalseKeyword,
	"finally":     ast.KindFinallyKeyword,
	"for":         ast.KindForKeyword,
	"from":        ast.KindFromKeyword,
	"function":    ast.KindFunctionKeyword,
	"get":         ast.KindGetKeyword,
	"if":          ast.KindIfKeyword,
	"immediate":   ast.KindImmediateKeyword,
	"implements":  ast.KindImplementsKeyword,
	"import":      ast.KindImportKeyword,
	"in":          ast.KindInKeyword,
	"infer":       ast.KindInferKeyword,
	"instanceof":  ast.KindInstanceOfKeyword,
	"interface":   ast.KindInterfaceKeyword,
	"intrinsic":   ast.KindIntrinsicKeyword,
	"is":          ast.KindIsKeyword,
	"keyof":       ast.KindKeyOfKeyword,
	"let":         ast.KindLetKeyword,
	"module":      ast.KindModuleKeyword,
	"namespace":   ast.KindNamespaceKeyword,
	"never":       ast.KindNeverKeyword,
	"new":         ast.KindNewKeyword,
	"null":        ast.KindNullKeyword,
	"number":      ast.KindNumberKeyword,
	"object":      ast.KindObjectKeyword,
	"package":     ast.KindPackageKeyword,
	"private":     ast.KindPrivateKeyword,
	"protected":   ast.KindProtectedKeyword,
	"public":      ast.KindPublicKeyword,
	"override":    ast.KindOverrideKeyword,
	"out":         ast.KindOutKeyword,
	"readonly":    ast.KindReadonlyKeyword,
	"require":     ast.KindRequireKeyword,
	"global":      ast.KindGlobalKeyword,
	"return":      ast.KindReturnKeyword,
	"satisfies":   ast.KindSatisfiesKeyword,
	"set":         ast.KindSetKeyword,
	"static":      ast.KindStaticKeyword,
	"string":      ast.KindStringKeyword,
	"super":       ast.KindSuperKeyword,
	"switch":      ast.KindSwitchKeyword,
	"symbol":      ast.KindSymbolKeyword,
	"this":        ast.KindThisKeyword,
	"throw":       ast.KindThrowKeyword,
	"true":        ast.KindTrueKeyword,
	"try":         ast.KindTryKeyword,
	"type":        ast.KindTypeKeyword,
	"typeof":      ast.KindTypeOfKeyword,
	"undefined":   ast.KindUndefinedKeyword,
	"unique":      ast.KindUniqueKeyword,
	"unknown":     ast.KindUnknownKeyword,
	"using":       ast.KindUsingKeyword,
	"var":         ast.KindVarKeyword,
	"void":        ast.KindVoidKeyword,
	"while":       ast.KindWhileKeyword,
	"with":        ast.KindWithKeyword,
	"yield":       ast.KindYieldKeyword,
	"async":       ast.KindAsyncKeyword,
	"await":       ast.KindAwaitKeyword,
	"of":          ast.KindOfKeyword,
}

var textToToken = func() map[string]ast.Kind {
	m := map[string]ast.Kind{
		"{":    ast.KindOpenBraceToken,
		"}":    ast.KindCloseBraceToken,
		"(":    ast.KindOpenParenToken,
		")":    ast.KindCloseParenToken,
		"[":    ast.KindOpenBracketToken,
		"]":    ast.KindCloseBracketToken,
		".":    ast.KindDotToken,
		"...":  ast.KindDotDotDotToken,
		";":    ast.KindSemicolonToken,
		",":    ast.KindCommaToken,
		"<":    ast.KindLessThanToken,
		">":    ast.KindGreaterThanToken,
		"<=":   ast.KindLessThanEqualsToken,
		">=":   ast.KindGreaterThanEqualsToken,
		"==":   ast.KindEqualsEqualsToken,
		"!=":   ast.KindExclamationEqualsToken,
		"===":  ast.KindEqualsEqualsEqualsToken,
		"!==":  ast.KindExclamationEqualsEqualsToken,
		"=>":   ast.KindEqualsGreaterThanToken,
		"+":    ast.KindPlusToken,
		"-":    ast.KindMinusToken,
		"**":   ast.KindAsteriskAsteriskToken,
		"*":    ast.KindAsteriskToken,
		"/":    ast.KindSlashToken,
		"%":    ast.KindPercentToken,
		"++":   ast.KindPlusPlusToken,
		"--":   ast.KindMinusMinusToken,
		"<<":   ast.KindLessThanLessThanToken,
		"</":   ast.KindLessThanSlashToken,
		">>":   ast.KindGreaterThanGreaterThanToken,
		">>>":  ast.KindGreaterThanGreaterThanGreaterThanToken,
		"&":    ast.KindAmpersandToken,
		"|":    ast.KindBarToken,
		"^":    ast.KindCaretToken,
		"!":    ast.KindExclamationToken,
		"~":    ast.KindTildeToken,
		"&&":   ast.KindAmpersandAmpersandToken,
		"||":   ast.KindBarBarToken,
		"?":    ast.KindQuestionToken,
		"??":   ast.KindQuestionQuestionToken,
		"?.":   ast.KindQuestionDotToken,
		":":    ast.KindColonToken,
		"=":    ast.KindEqualsToken,
		"+=":   ast.KindPlusEqualsToken,
		"-=":   ast.KindMinusEqualsToken,
		"*=":   ast.KindAsteriskEqualsToken,
		"**=":  ast.KindAsteriskAsteriskEqualsToken,
		"/=":   ast.KindSlashEqualsToken,
		"%=":   ast.KindPercentEqualsToken,
		"<<=":  ast.KindLessThanLessThanEqualsToken,
		">>=":  ast.KindGreaterThanGreaterThanEqualsToken,
		">>>=": ast.KindGreaterThanGreaterThanGreaterThanEqualsToken,
		"&=":   ast.KindAmpersandEqualsToken,
		"|=":   ast.KindBarEqualsToken,
		"^=":   ast.KindCaretEqualsToken,
		"||=":  ast.KindBarBarEqualsToken,
		"&&=":  ast.KindAmpersandAmpersandEqualsToken,
		"??=":  ast.KindQuestionQuestionEqualsToken,
		"@":    ast.KindAtToken,
		"#":    ast.KindHashToken,
		"`":    ast.KindBacktickToken,
	}
	maps.Copy(m, textToKeyword)
	return m
}()

// Generated by scripts/regenerate-unicode-identifier-parts.mjs on node v22.1.0 with unicode 15.1
// based on http://www.unicode.org/reports/tr31/ and https://www.ecma-international.org/ecma-262/6.0/#sec-names-and-keywords
// unicodeESNextIdentifierStart corresponds to the ID_Start and Other_ID_Start property, and
// unicodeESNextIdentifierPart corresponds to ID_Continue, Other_ID_Continue, plus ID_Start and Other_ID_Start
var (
	unicodeESNextIdentifierStart = []rune{65, 90, 97, 122, 170, 170, 181, 181, 186, 186, 192, 214, 216, 246, 248, 705, 710, 721, 736, 740, 748, 748, 750, 750, 880, 884, 886, 887, 890, 893, 895, 895, 902, 902, 904, 906, 908, 908, 910, 929, 931, 1013, 1015, 1153, 1162, 1327, 1329, 1366, 1369, 1369, 1376, 1416, 1488, 1514, 1519, 1522, 1568, 1610, 1646, 1647, 1649, 1747, 1749, 1749, 1765, 1766, 1774, 1775, 1786, 1788, 1791, 1791, 1808, 1808, 1810, 1839, 1869, 1957, 1969, 1969, 1994, 2026, 2036, 2037, 2042, 2042, 2048, 2069, 2074, 2074, 2084, 2084, 2088, 2088, 2112, 2136, 2144, 2154, 2160, 2183, 2185, 2190, 2208, 2249, 2308, 2361, 2365, 2365, 2384, 2384, 2392, 2401, 2417, 2432, 2437, 2444, 2447, 2448, 2451, 2472, 2474, 2480, 2482, 2482, 2486, 2489, 2493, 2493, 2510, 2510, 2524, 2525, 2527, 2529, 2544, 2545, 2556, 2556, 2565, 2570, 2575, 2576, 2579, 2600, 2602, 2608, 2610, 2611, 2613, 2614, 2616, 2617, 2649, 2652, 2654, 2654, 2674, 2676, 2693, 2701, 2703, 2705, 2707, 2728, 2730, 2736, 2738, 2739, 2741, 2745, 2749, 2749, 2768, 2768, 2784, 2785, 2809, 2809, 2821, 2828, 2831, 2832, 2835, 2856, 2858, 2864, 2866, 2867, 2869, 2873, 2877, 2877, 2908, 2909, 2911, 2913, 2929, 2929, 2947, 2947, 2949, 2954, 2958, 2960, 2962, 2965, 2969, 2970, 2972, 2972, 2974, 2975, 2979, 2980, 2984, 2986, 2990, 3001, 3024, 3024, 3077, 3084, 3086, 3088, 3090, 3112, 3114, 3129, 3133, 3133, 3160, 3162, 3165, 3165, 3168, 3169, 3200, 3200, 3205, 3212, 3214, 3216, 3218, 3240, 3242, 3251, 3253, 3257, 3261, 3261, 3293, 3294, 3296, 3297, 3313, 3314, 3332, 3340, 3342, 3344, 3346, 3386, 3389, 3389, 3406, 3406, 3412, 3414, 3423, 3425, 3450, 3455, 3461, 3478, 3482, 3505, 3507, 3515, 3517, 3517, 3520, 3526, 3585, 3632, 3634, 3635, 3648, 3654, 3713, 3714, 3716, 3716, 3718, 3722, 3724, 3747, 3749, 3749, 3751, 3760, 3762, 3763, 3773, 3773, 3776, 3780, 3782, 3782, 3804, 3807, 3840, 3840, 3904, 3911, 3913, 3948, 3976, 3980, 4096, 4138, 4159, 4159, 4176, 4181, 4186, 4189, 4193, 4193, 4197, 4198, 4206, 4208, 4213, 4225, 4238, 4238, 4256, 4293, 4295, 4295, 4301, 4301, 4304, 4346, 4348, 4680, 4682, 4685, 4688, 4694, 4696, 4696, 4698, 4701, 4704, 4744, 4746, 4749, 4752, 4784, 4786, 4789, 4792, 4798, 4800, 4800, 4802, 4805, 4808, 4822, 4824, 4880, 4882, 4885, 4888, 4954, 4992, 5007, 5024, 5109, 5112, 5117, 5121, 5740, 5743, 5759, 5761, 5786, 5792, 5866, 5870, 5880, 5888, 5905, 5919, 5937, 5952, 5969, 5984, 5996, 5998, 6000, 6016, 6067, 6103, 6103, 6108, 6108, 6176, 6264, 6272, 6312, 6314, 6314, 6320, 6389, 6400, 6430, 6480, 6509, 6512, 6516, 6528, 6571, 6576, 6601, 6656, 6678, 6688, 6740, 6823, 6823, 6917, 6963, 6981, 6988, 7043, 7072, 7086, 7087, 7098, 7141, 7168, 7203, 7245, 7247, 7258, 7293, 7296, 7304, 7312, 7354, 7357, 7359, 7401, 7404, 7406, 7411, 7413, 7414, 7418, 7418, 7424, 7615, 7680, 7957, 7960, 7965, 7968, 8005, 8008, 8013, 8016, 8023, 8025, 8025, 8027, 8027, 8029, 8029, 8031, 8061, 8064, 8116, 8118, 8124, 8126, 8126, 8130, 8132, 8134, 8140, 8144, 8147, 8150, 8155, 8160, 8172, 8178, 8180, 8182, 8188, 8305, 8305, 8319, 8319, 8336, 8348, 8450, 8450, 8455, 8455, 8458, 8467, 8469, 8469, 8472, 8477, 8484, 8484, 8486, 8486, 8488, 8488, 8490, 8505, 8508, 8511, 8517, 8521, 8526, 8526, 8544, 8584, 11264, 11492, 11499, 11502, 11506, 11507, 11520, 11557, 11559, 11559, 11565, 11565, 11568, 11623, 11631, 11631, 11648, 11670, 11680, 11686, 11688, 11694, 11696, 11702, 11704, 11710, 11712, 11718, 11720, 11726, 11728, 11734, 11736, 11742, 12293, 12295, 12321, 12329, 12337, 12341, 12344, 12348, 12353, 12438, 12443, 12447, 12449, 12538, 12540, 12543, 12549, 12591, 12593, 12686, 12704, 12735, 12784, 12799, 13312, 19903, 19968, 42124, 42192, 42237, 42240, 42508, 42512, 42527, 42538, 42539, 42560, 42606, 42623, 42653, 42656, 42735, 42775, 42783, 42786, 42888, 42891, 42954, 42960, 42961, 42963, 42963, 42965, 42969, 42994, 43009, 43011, 43013, 43015, 43018, 43020, 43042, 43072, 43123, 43138, 43187, 43250, 43255, 43259, 43259, 43261, 43262, 43274, 43301, 43312, 43334, 43360, 43388, 43396, 43442, 43471, 43471, 43488, 43492, 43494, 43503, 43514, 43518, 43520, 43560, 43584, 43586, 43588, 43595, 43616, 43638, 43642, 43642, 43646, 43695, 43697, 43697, 43701, 43702, 43705, 43709, 43712, 43712, 43714, 43714, 43739, 43741, 43744, 43754, 43762, 43764, 43777, 43782, 43785, 43790, 43793, 43798, 43808, 43814, 43816, 43822, 43824, 43866, 43868, 43881, 43888, 44002, 44032, 55203, 55216, 55238, 55243, 55291, 63744, 64109, 64112, 64217, 64256, 64262, 64275, 64279, 64285, 64285, 64287, 64296, 64298, 64310, 64312, 64316, 64318, 64318, 64320, 64321, 64323, 64324, 64326, 64433, 64467, 64829, 64848, 64911, 64914, 64967, 65008, 65019, 65136, 65140, 65142, 65276, 65313, 65338, 65345, 65370, 65382, 65470, 65474, 65479, 65482, 65487, 65490, 65495, 65498, 65500, 65536, 65547, 65549, 65574, 65576, 65594, 65596, 65597, 65599, 65613, 65616, 65629, 65664, 65786, 65856, 65908, 66176, 66204, 66208, 66256, 66304, 66335, 66349, 66378, 66384, 66421, 66432, 66461, 66464, 66499, 66504, 66511, 66513, 66517, 66560, 66717, 66736, 66771, 66776, 66811, 66816, 66855, 66864, 66915, 66928, 66938, 66940, 66954, 66956, 66962, 66964, 66965, 66967, 66977, 66979, 66993, 66995, 67001, 67003, 67004, 67072, 67382, 67392, 67413, 67424, 67431, 67456, 67461, 67463, 67504, 67506, 67514, 67584, 67589, 67592, 67592, 67594, 67637, 67639, 67640, 67644, 67644, 67647, 67669, 67680, 67702, 67712, 67742, 67808, 67826, 67828, 67829, 67840, 67861, 67872, 67897, 67968, 68023, 68030, 68031, 68096, 68096, 68112, 68115, 68117, 68119, 68121, 68149, 68192, 68220, 68224, 68252, 68288, 68295, 68297, 68324, 68352, 68405, 68416, 68437, 68448, 68466, 68480, 68497, 68608, 68680, 68736, 68786, 68800, 68850, 68864, 68899, 69248, 69289, 69296, 69297, 69376, 69404, 69415, 69415, 69424, 69445, 69488, 69505, 69552, 69572, 69600, 69622, 69635, 69687, 69745, 69746, 69749, 69749, 69763, 69807, 69840, 69864, 69891, 69926, 69956, 69956, 69959, 69959, 69968, 70002, 70006, 70006, 70019, 70066, 70081, 70084, 70106, 70106, 70108, 70108, 70144, 70161, 70163, 70187, 70207, 70208, 70272, 70278, 70280, 70280, 70282, 70285, 70287, 70301, 70303, 70312, 70320, 70366, 70405, 70412, 70415, 70416, 70419, 70440, 70442, 70448, 70450, 70451, 70453, 70457, 70461, 70461, 70480, 70480, 70493, 70497, 70656, 70708, 70727, 70730, 70751, 70753, 70784, 70831, 70852, 70853, 70855, 70855, 71040, 71086, 71128, 71131, 71168, 71215, 71236, 71236, 71296, 71338, 71352, 71352, 71424, 71450, 71488, 71494, 71680, 71723, 71840, 71903, 71935, 71942, 71945, 71945, 71948, 71955, 71957, 71958, 71960, 71983, 71999, 71999, 72001, 72001, 72096, 72103, 72106, 72144, 72161, 72161, 72163, 72163, 72192, 72192, 72203, 72242, 72250, 72250, 72272, 72272, 72284, 72329, 72349, 72349, 72368, 72440, 72704, 72712, 72714, 72750, 72768, 72768, 72818, 72847, 72960, 72966, 72968, 72969, 72971, 73008, 73030, 73030, 73056, 73061, 73063, 73064, 73066, 73097, 73112, 73112, 73440, 73458, 73474, 73474, 73476, 73488, 73490, 73523, 73648, 73648, 73728, 74649, 74752, 74862, 74880, 75075, 77712, 77808, 77824, 78895, 78913, 78918, 82944, 83526, 92160, 92728, 92736, 92766, 92784, 92862, 92880, 92909, 92928, 92975, 92992, 92995, 93027, 93047, 93053, 93071, 93760, 93823, 93952, 94026, 94032, 94032, 94099, 94111, 94176, 94177, 94179, 94179, 94208, 100343, 100352, 101589, 101632, 101640, 110576, 110579, 110581, 110587, 110589, 110590, 110592, 110882, 110898, 110898, 110928, 110930, 110933, 110933, 110948, 110951, 110960, 111355, 113664, 113770, 113776, 113788, 113792, 113800, 113808, 113817, 119808, 119892, 119894, 119964, 119966, 119967, 119970, 119970, 119973, 119974, 119977, 119980, 119982, 119993, 119995, 119995, 119997, 120003, 120005, 120069, 120071, 120074, 120077, 120084, 120086, 120092, 120094, 120121, 120123, 120126, 120128, 120132, 120134, 120134, 120138, 120144, 120146, 120485, 120488, 120512, 120514, 120538, 120540, 120570, 120572, 120596, 120598, 120628, 120630, 120654, 120656, 120686, 120688, 120712, 120714, 120744, 120746, 120770, 120772, 120779, 122624, 122654, 122661, 122666, 122928, 122989, 123136, 123180, 123191, 123197, 123214, 123214, 123536, 123565, 123584, 123627, 124112, 124139, 124896, 124902, 124904, 124907, 124909, 124910, 124912, 124926, 124928, 125124, 125184, 125251, 125259, 125259, 126464, 126467, 126469, 126495, 126497, 126498, 126500, 126500, 126503, 126503, 126505, 126514, 126516, 126519, 126521, 126521, 126523, 126523, 126530, 126530, 126535, 126535, 126537, 126537, 126539, 126539, 126541, 126543, 126545, 126546, 126548, 126548, 126551, 126551, 126553, 126553, 126555, 126555, 126557, 126557, 126559, 126559, 126561, 126562, 126564, 126564, 126567, 126570, 126572, 126578, 126580, 126583, 126585, 126588, 126590, 126590, 126592, 126601, 126603, 126619, 126625, 126627, 126629, 126633, 126635, 126651, 131072, 173791, 173824, 177977, 177984, 178205, 178208, 183969, 183984, 191456, 191472, 192093, 194560, 195101, 196608, 201546, 201552, 205743}
	unicodeESNextIdentifierPart  = []rune{48, 57, 65, 90, 95, 95, 97, 122, 170, 170, 181, 181, 183, 183, 186, 186, 192, 214, 216, 246, 248, 705, 710, 721, 736, 740, 748, 748, 750, 750, 768, 884, 886, 887, 890, 893, 895, 895, 902, 906, 908, 908, 910, 929, 931, 1013, 1015, 1153, 1155, 1159, 1162, 1327, 1329, 1366, 1369, 1369, 1376, 1416, 1425, 1469, 1471, 1471, 1473, 1474, 1476, 1477, 1479, 1479, 1488, 1514, 1519, 1522, 1552, 1562, 1568, 1641, 1646, 1747, 1749, 1756, 1759, 1768, 1770, 1788, 1791, 1791, 1808, 1866, 1869, 1969, 1984, 2037, 2042, 2042, 2045, 2045, 2048, 2093, 2112, 2139, 2144, 2154, 2160, 2183, 2185, 2190, 2200, 2273, 2275, 2403, 2406, 2415, 2417, 2435, 2437, 2444, 2447, 2448, 2451, 2472, 2474, 2480, 2482, 2482, 2486, 2489, 2492, 2500, 2503, 2504, 2507, 2510, 2519, 2519, 2524, 2525, 2527, 2531, 2534, 2545, 2556, 2556, 2558, 2558, 2561, 2563, 2565, 2570, 2575, 2576, 2579, 2600, 2602, 2608, 2610, 2611, 2613, 2614, 2616, 2617, 2620, 2620, 2622, 2626, 2631, 2632, 2635, 2637, 2641, 2641, 2649, 2652, 2654, 2654, 2662, 2677, 2689, 2691, 2693, 2701, 2703, 2705, 2707, 2728, 2730, 2736, 2738, 2739, 2741, 2745, 2748, 2757, 2759, 2761, 2763, 2765, 2768, 2768, 2784, 2787, 2790, 2799, 2809, 2815, 2817, 2819, 2821, 2828, 2831, 2832, 2835, 2856, 2858, 2864, 2866, 2867, 2869, 2873, 2876, 2884, 2887, 2888, 2891, 2893, 2901, 2903, 2908, 2909, 2911, 2915, 2918, 2927, 2929, 2929, 2946, 2947, 2949, 2954, 2958, 2960, 2962, 2965, 2969, 2970, 2972, 2972, 2974, 2975, 2979, 2980, 2984, 2986, 2990, 3001, 3006, 3010, 3014, 3016, 3018, 3021, 3024, 3024, 3031, 3031, 3046, 3055, 3072, 3084, 3086, 3088, 3090, 3112, 3114, 3129, 3132, 3140, 3142, 3144, 3146, 3149, 3157, 3158, 3160, 3162, 3165, 3165, 3168, 3171, 3174, 3183, 3200, 3203, 3205, 3212, 3214, 3216, 3218, 3240, 3242, 3251, 3253, 3257, 3260, 3268, 3270, 3272, 3274, 3277, 3285, 3286, 3293, 3294, 3296, 3299, 3302, 3311, 3313, 3315, 3328, 3340, 3342, 3344, 3346, 3396, 3398, 3400, 3402, 3406, 3412, 3415, 3423, 3427, 3430, 3439, 3450, 3455, 3457, 3459, 3461, 3478, 3482, 3505, 3507, 3515, 3517, 3517, 3520, 3526, 3530, 3530, 3535, 3540, 3542, 3542, 3544, 3551, 3558, 3567, 3570, 3571, 3585, 3642, 3648, 3662, 3664, 3673, 3713, 3714, 3716, 3716, 3718, 3722, 3724, 3747, 3749, 3749, 3751, 3773, 3776, 3780, 3782, 3782, 3784, 3790, 3792, 3801, 3804, 3807, 3840, 3840, 3864, 3865, 3872, 3881, 3893, 3893, 3895, 3895, 3897, 3897, 3902, 3911, 3913, 3948, 3953, 3972, 3974, 3991, 3993, 4028, 4038, 4038, 4096, 4169, 4176, 4253, 4256, 4293, 4295, 4295, 4301, 4301, 4304, 4346, 4348, 4680, 4682, 4685, 4688, 4694, 4696, 4696, 4698, 4701, 4704, 4744, 4746, 4749, 4752, 4784, 4786, 4789, 4792, 4798, 4800, 4800, 4802, 4805, 4808, 4822, 4824, 4880, 4882, 4885, 4888, 4954, 4957, 4959, 4969, 4977, 4992, 5007, 5024, 5109, 5112, 5117, 5121, 5740, 5743, 5759, 5761, 5786, 5792, 5866, 5870, 5880, 5888, 5909, 5919, 5940, 5952, 5971, 5984, 5996, 5998, 6000, 6002, 6003, 6016, 6099, 6103, 6103, 6108, 6109, 6112, 6121, 6155, 6157, 6159, 6169, 6176, 6264, 6272, 6314, 6320, 6389, 6400, 6430, 6432, 6443, 6448, 6459, 6470, 6509, 6512, 6516, 6528, 6571, 6576, 6601, 6608, 6618, 6656, 6683, 6688, 6750, 6752, 6780, 6783, 6793, 6800, 6809, 6823, 6823, 6832, 6845, 6847, 6862, 6912, 6988, 6992, 7001, 7019, 7027, 7040, 7155, 7168, 7223, 7232, 7241, 7245, 7293, 7296, 7304, 7312, 7354, 7357, 7359, 7376, 7378, 7380, 7418, 7424, 7957, 7960, 7965, 7968, 8005, 8008, 8013, 8016, 8023, 8025, 8025, 8027, 8027, 8029, 8029, 8031, 8061, 8064, 8116, 8118, 8124, 8126, 8126, 8130, 8132, 8134, 8140, 8144, 8147, 8150, 8155, 8160, 8172, 8178, 8180, 8182, 8188, 8204, 8205, 8255, 8256, 8276, 8276, 8305, 8305, 8319, 8319, 8336, 8348, 8400, 8412, 8417, 8417, 8421, 8432, 8450, 8450, 8455, 8455, 8458, 8467, 8469, 8469, 8472, 8477, 8484, 8484, 8486, 8486, 8488, 8488, 8490, 8505, 8508, 8511, 8517, 8521, 8526, 8526, 8544, 8584, 11264, 11492, 11499, 11507, 11520, 11557, 11559, 11559, 11565, 11565, 11568, 11623, 11631, 11631, 11647, 11670, 11680, 11686, 11688, 11694, 11696, 11702, 11704, 11710, 11712, 11718, 11720, 11726, 11728, 11734, 11736, 11742, 11744, 11775, 12293, 12295, 12321, 12335, 12337, 12341, 12344, 12348, 12353, 12438, 12441, 12447, 12449, 12543, 12549, 12591, 12593, 12686, 12704, 12735, 12784, 12799, 13312, 19903, 19968, 42124, 42192, 42237, 42240, 42508, 42512, 42539, 42560, 42607, 42612, 42621, 42623, 42737, 42775, 42783, 42786, 42888, 42891, 42954, 42960, 42961, 42963, 42963, 42965, 42969, 42994, 43047, 43052, 43052, 43072, 43123, 43136, 43205, 43216, 43225, 43232, 43255, 43259, 43259, 43261, 43309, 43312, 43347, 43360, 43388, 43392, 43456, 43471, 43481, 43488, 43518, 43520, 43574, 43584, 43597, 43600, 43609, 43616, 43638, 43642, 43714, 43739, 43741, 43744, 43759, 43762, 43766, 43777, 43782, 43785, 43790, 43793, 43798, 43808, 43814, 43816, 43822, 43824, 43866, 43868, 43881, 43888, 44010, 44012, 44013, 44016, 44025, 44032, 55203, 55216, 55238, 55243, 55291, 63744, 64109, 64112, 64217, 64256, 64262, 64275, 64279, 64285, 64296, 64298, 64310, 64312, 64316, 64318, 64318, 64320, 64321, 64323, 64324, 64326, 64433, 64467, 64829, 64848, 64911, 64914, 64967, 65008, 65019, 65024, 65039, 65056, 65071, 65075, 65076, 65101, 65103, 65136, 65140, 65142, 65276, 65296, 65305, 65313, 65338, 65343, 65343, 65345, 65370, 65381, 65470, 65474, 65479, 65482, 65487, 65490, 65495, 65498, 65500, 65536, 65547, 65549, 65574, 65576, 65594, 65596, 65597, 65599, 65613, 65616, 65629, 65664, 65786, 65856, 65908, 66045, 66045, 66176, 66204, 66208, 66256, 66272, 66272, 66304, 66335, 66349, 66378, 66384, 66426, 66432, 66461, 66464, 66499, 66504, 66511, 66513, 66517, 66560, 66717, 66720, 66729, 66736, 66771, 66776, 66811, 66816, 66855, 66864, 66915, 66928, 66938, 66940, 66954, 66956, 66962, 66964, 66965, 66967, 66977, 66979, 66993, 66995, 67001, 67003, 67004, 67072, 67382, 67392, 67413, 67424, 67431, 67456, 67461, 67463, 67504, 67506, 67514, 67584, 67589, 67592, 67592, 67594, 67637, 67639, 67640, 67644, 67644, 67647, 67669, 67680, 67702, 67712, 67742, 67808, 67826, 67828, 67829, 67840, 67861, 67872, 67897, 67968, 68023, 68030, 68031, 68096, 68099, 68101, 68102, 68108, 68115, 68117, 68119, 68121, 68149, 68152, 68154, 68159, 68159, 68192, 68220, 68224, 68252, 68288, 68295, 68297, 68326, 68352, 68405, 68416, 68437, 68448, 68466, 68480, 68497, 68608, 68680, 68736, 68786, 68800, 68850, 68864, 68903, 68912, 68921, 69248, 69289, 69291, 69292, 69296, 69297, 69373, 69404, 69415, 69415, 69424, 69456, 69488, 69509, 69552, 69572, 69600, 69622, 69632, 69702, 69734, 69749, 69759, 69818, 69826, 69826, 69840, 69864, 69872, 69881, 69888, 69940, 69942, 69951, 69956, 69959, 69968, 70003, 70006, 70006, 70016, 70084, 70089, 70092, 70094, 70106, 70108, 70108, 70144, 70161, 70163, 70199, 70206, 70209, 70272, 70278, 70280, 70280, 70282, 70285, 70287, 70301, 70303, 70312, 70320, 70378, 70384, 70393, 70400, 70403, 70405, 70412, 70415, 70416, 70419, 70440, 70442, 70448, 70450, 70451, 70453, 70457, 70459, 70468, 70471, 70472, 70475, 70477, 70480, 70480, 70487, 70487, 70493, 70499, 70502, 70508, 70512, 70516, 70656, 70730, 70736, 70745, 70750, 70753, 70784, 70853, 70855, 70855, 70864, 70873, 71040, 71093, 71096, 71104, 71128, 71133, 71168, 71232, 71236, 71236, 71248, 71257, 71296, 71352, 71360, 71369, 71424, 71450, 71453, 71467, 71472, 71481, 71488, 71494, 71680, 71738, 71840, 71913, 71935, 71942, 71945, 71945, 71948, 71955, 71957, 71958, 71960, 71989, 71991, 71992, 71995, 72003, 72016, 72025, 72096, 72103, 72106, 72151, 72154, 72161, 72163, 72164, 72192, 72254, 72263, 72263, 72272, 72345, 72349, 72349, 72368, 72440, 72704, 72712, 72714, 72758, 72760, 72768, 72784, 72793, 72818, 72847, 72850, 72871, 72873, 72886, 72960, 72966, 72968, 72969, 72971, 73014, 73018, 73018, 73020, 73021, 73023, 73031, 73040, 73049, 73056, 73061, 73063, 73064, 73066, 73102, 73104, 73105, 73107, 73112, 73120, 73129, 73440, 73462, 73472, 73488, 73490, 73530, 73534, 73538, 73552, 73561, 73648, 73648, 73728, 74649, 74752, 74862, 74880, 75075, 77712, 77808, 77824, 78895, 78912, 78933, 82944, 83526, 92160, 92728, 92736, 92766, 92768, 92777, 92784, 92862, 92864, 92873, 92880, 92909, 92912, 92916, 92928, 92982, 92992, 92995, 93008, 93017, 93027, 93047, 93053, 93071, 93760, 93823, 93952, 94026, 94031, 94087, 94095, 94111, 94176, 94177, 94179, 94180, 94192, 94193, 94208, 100343, 100352, 101589, 101632, 101640, 110576, 110579, 110581, 110587, 110589, 110590, 110592, 110882, 110898, 110898, 110928, 110930, 110933, 110933, 110948, 110951, 110960, 111355, 113664, 113770, 113776, 113788, 113792, 113800, 113808, 113817, 113821, 113822, 118528, 118573, 118576, 118598, 119141, 119145, 119149, 119154, 119163, 119170, 119173, 119179, 119210, 119213, 119362, 119364, 119808, 119892, 119894, 119964, 119966, 119967, 119970, 119970, 119973, 119974, 119977, 119980, 119982, 119993, 119995, 119995, 119997, 120003, 120005, 120069, 120071, 120074, 120077, 120084, 120086, 120092, 120094, 120121, 120123, 120126, 120128, 120132, 120134, 120134, 120138, 120144, 120146, 120485, 120488, 120512, 120514, 120538, 120540, 120570, 120572, 120596, 120598, 120628, 120630, 120654, 120656, 120686, 120688, 120712, 120714, 120744, 120746, 120770, 120772, 120779, 120782, 120831, 121344, 121398, 121403, 121452, 121461, 121461, 121476, 121476, 121499, 121503, 121505, 121519, 122624, 122654, 122661, 122666, 122880, 122886, 122888, 122904, 122907, 122913, 122915, 122916, 122918, 122922, 122928, 122989, 123023, 123023, 123136, 123180, 123184, 123197, 123200, 123209, 123214, 123214, 123536, 123566, 123584, 123641, 124112, 124153, 124896, 124902, 124904, 124907, 124909, 124910, 124912, 124926, 124928, 125124, 125136, 125142, 125184, 125259, 125264, 125273, 126464, 126467, 126469, 126495, 126497, 126498, 126500, 126500, 126503, 126503, 126505, 126514, 126516, 126519, 126521, 126521, 126523, 126523, 126530, 126530, 126535, 126535, 126537, 126537, 126539, 126539, 126541, 126543, 126545, 126546, 126548, 126548, 126551, 126551, 126553, 126553, 126555, 126555, 126557, 126557, 126559, 126559, 126561, 126562, 126564, 126564, 126567, 126570, 126572, 126578, 126580, 126583, 126585, 126588, 126590, 126590, 126592, 126601, 126603, 126619, 126625, 126627, 126629, 126633, 126635, 126651, 130032, 130041, 131072, 173791, 173824, 177977, 177984, 178205, 178208, 183969, 183984, 191456, 191472, 192093, 194560, 195101, 196608, 201546, 201552, 205743, 917760, 917999}
)

type ScannerState struct {
	pos                       int            // Current position in text (and ending position of current token)
	fullStartPos              int            // Starting position of current token including preceding whitespace
	tokenStart                int            // Starting position of non-whitespace part of current token
	token                     ast.Kind       // Kind of current token
	tokenValue                string         // Parsed value of current token
	tokenFlags                ast.TokenFlags // Flags for current token
	commentDirectives         []ast.CommentDirective
	skipJSDocLeadingAsterisks int // Leading asterisks to skip when scanning types inside JSDoc. Should be 0 outside JSDoc
}

type Scanner struct {
	text             string
	languageVariant  core.LanguageVariant
	onError          ErrorCallback
	skipTrivia       bool
	JSDocParsingMode ast.JSDocParsingMode
	scriptKind       core.ScriptKind
	ScannerState

	numberCache    map[string]string
	hexNumberCache map[string]string
	hexDigitCache  map[string]string
}

func defaultScanner() Scanner {
	// Using a function rather than a global is intentional; this function is
	// inlined as pure code (zeroing + moves), whereas a global requires write
	// barriers since the memory is mutable.
	return Scanner{skipTrivia: true}
}

func NewScanner() *Scanner {
	s := defaultScanner()
	return &s
}

func (s *Scanner) Reset() {
	numberCache := cleared(s.numberCache)
	hexNumberCache := cleared(s.hexNumberCache)
	hexDigitCache := cleared(s.hexDigitCache)
	*s = defaultScanner()
	s.numberCache = numberCache
	s.hexNumberCache = hexNumberCache
	s.hexDigitCache = hexDigitCache
}

func cleared[M ~map[K]V, K comparable, V any](m M) M {
	clear(m)
	return m
}

func (s *Scanner) Text() string {
	return s.text
}

func (s *Scanner) Token() ast.Kind {
	return s.token
}

func (s *Scanner) TokenFlags() ast.TokenFlags {
	return s.tokenFlags
}

func (s *Scanner) TokenFullStart() int {
	return s.fullStartPos
}

func (s *Scanner) TokenStart() int {
	return s.tokenStart
}

func (s *Scanner) TokenEnd() int {
	return s.pos
}

func (s *Scanner) TokenText() string {
	return s.text[s.tokenStart:s.pos]
}

func (s *Scanner) TokenValue() string {
	return s.tokenValue
}

func (s *Scanner) TokenRange() core.TextRange {
	return core.NewTextRange(s.tokenStart, s.pos)
}

func (s *Scanner) CommentDirectives() []ast.CommentDirective {
	return s.commentDirectives
}

func (s *Scanner) Mark() ScannerState {
	return s.ScannerState
}

func (s *Scanner) Rewind(state ScannerState) {
	s.ScannerState = state
}

func (s *Scanner) ResetPos(pos int) {
	if pos < 0 {
		panic("Cannot reset token state to negative position")
	}
	s.pos = pos
	s.fullStartPos = pos
	s.tokenStart = pos
}

func (s *Scanner) ResetTokenState(pos int) {
	s.ResetPos(pos)
	s.token = ast.KindUnknown
	s.tokenValue = ""
	s.tokenFlags = ast.TokenFlagsNone
}

func (scanner *Scanner) SetSkipJSDocLeadingAsterisks(skip bool) {
	if skip {
		scanner.skipJSDocLeadingAsterisks += 1
	} else {
		scanner.skipJSDocLeadingAsterisks += -1
	}
}

func (scanner *Scanner) SetSkipTrivia(skip bool) {
	scanner.skipTrivia = skip
}

func (s *Scanner) HasUnicodeEscape() bool {
	return s.tokenFlags&ast.TokenFlagsUnicodeEscape != 0
}

func (s *Scanner) HasExtendedUnicodeEscape() bool {
	return s.tokenFlags&ast.TokenFlagsExtendedUnicodeEscape != 0
}

func (s *Scanner) HasPrecedingLineBreak() bool {
	return s.tokenFlags&ast.TokenFlagsPrecedingLineBreak != 0
}

func (s *Scanner) HasPrecedingJSDocComment() bool {
	return s.tokenFlags&ast.TokenFlagsPrecedingJSDocComment != 0
}

func (s *Scanner) HasPrecedingJSDocLeadingAsterisks() bool {
	return s.tokenFlags&ast.TokenFlagsPrecedingJSDocLeadingAsterisks != 0
}

func (s *Scanner) SetText(text string) {
	s.text = text
	s.ScannerState = ScannerState{}
}

func (s *Scanner) SetOnError(errorCallback ErrorCallback) {
	s.onError = errorCallback
}

func (s *Scanner) SetScriptKind(scriptKind core.ScriptKind) {
	s.scriptKind = scriptKind
}

func (s *Scanner) SetJSDocParsingMode(kind ast.JSDocParsingMode) {
	s.JSDocParsingMode = kind
}

func (s *Scanner) SetLanguageVariant(languageVariant core.LanguageVariant) {
	s.languageVariant = languageVariant
}

func (s *Scanner) error(diagnostic *diagnostics.Message) {
	s.errorAt(diagnostic, s.pos, 0)
}

func (s *Scanner) errorAt(diagnostic *diagnostics.Message, pos int, length int, args ...any) {
	if s.onError != nil {
		s.onError(diagnostic, pos, length, args...)
	}
}

// NOTE: even though this returns a rune, it only decodes the current byte.
// It must be checked against utf8.RuneSelf to verify that a call to charAndSize
// is not needed.
func (s *Scanner) char() rune {
	if s.pos < len(s.text) {
		return rune(s.text[s.pos])
	}
	return -1
}

// NOTE: this returns a rune, but only decodes the byte at the offset.
func (s *Scanner) charAt(offset int) rune {
	if s.pos+offset < len(s.text) {
		return rune(s.text[s.pos+offset])
	}
	return -1
}

func (s *Scanner) charAndSize() (rune, int) {
	return utf8.DecodeRuneInString(s.text[s.pos:])
}

func (s *Scanner) shouldParseJSDoc() bool {
	switch s.JSDocParsingMode {
	case ast.JSDocParsingModeParseAll:
		return true
	case ast.JSDocParsingModeParseNone:
		return false
	}
	if s.scriptKind != core.ScriptKindTS && s.scriptKind != core.ScriptKindTSX {
		// If outside of TS, we need JSDoc to get any type info.
		return true
	}
	if s.JSDocParsingMode == ast.JSDocParsingModeParseForTypeInfo {
		// If we're in TS, but we don't need to produce reliable errors,
		// we don't need to parse to find @see or @link.
		return false
	}
	text := s.text[s.fullStartPos:s.pos]
	for {
		i := strings.IndexByte(text, '@')
		if i < 0 {
			break
		}
		text = text[i+1:]
		if strings.HasPrefix(text, "see") || strings.HasPrefix(text, "link") {
			return true
		}
	}
	return false
}

func (s *Scanner) Scan() ast.Kind {
	s.fullStartPos = s.pos
	s.tokenFlags = ast.TokenFlagsNone
	for {
		s.tokenStart = s.pos
		ch := s.char()

		switch ch {
		case '\t', '\v', '\f', ' ':
			s.pos++
			if s.skipTrivia {
				continue
			}
			for {
				ch, size := s.charAndSize()
				if !stringutil.IsWhiteSpaceSingleLine(ch) {
					break
				}
				s.pos += size
			}
			s.token = ast.KindWhitespaceTrivia
		case '\n', '\r':
			s.tokenFlags |= ast.TokenFlagsPrecedingLineBreak
			if s.skipTrivia {
				s.pos++
				continue
			}
			if ch == '\r' && s.charAt(1) == '\n' {
				s.pos += 2
			} else {
				s.pos++
			}
			s.token = ast.KindNewLineTrivia
		case '!':
			if s.charAt(1) == '=' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindExclamationEqualsEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindExclamationEqualsToken
				}
			} else {
				s.pos++
				s.token = ast.KindExclamationToken
			}
		case '"', '\'':
			s.tokenValue = s.scanString(false /*jsxAttributeString*/)
			s.token = ast.KindStringLiteral
		case '`':
			s.token = s.scanTemplateAndSetTokenValue(false /*shouldEmitInvalidEscapeError*/)
		case '%':
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindPercentEqualsToken
			} else {
				s.pos++
				s.token = ast.KindPercentToken
			}
		case '&':
			if s.charAt(1) == '&' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindAmpersandAmpersandEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindAmpersandAmpersandToken
				}
			} else if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindAmpersandEqualsToken
			} else {
				s.pos++
				s.token = ast.KindAmpersandToken
			}
		case '(':
			s.pos++
			s.token = ast.KindOpenParenToken
		case ')':
			s.pos++
			s.token = ast.KindCloseParenToken
		case '*':
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindAsteriskEqualsToken
			} else if s.charAt(1) == '*' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindAsteriskAsteriskEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindAsteriskAsteriskToken
				}
			} else {
				s.pos++
				if s.skipJSDocLeadingAsterisks != 0 &&
					(s.tokenFlags&ast.TokenFlagsPrecedingJSDocLeadingAsterisks) == 0 &&
					(s.tokenFlags&ast.TokenFlagsPrecedingLineBreak) != 0 {
					s.tokenFlags |= ast.TokenFlagsPrecedingJSDocLeadingAsterisks
					continue
				}
				s.token = ast.KindAsteriskToken
			}
		case '+':
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindPlusEqualsToken
			} else if s.charAt(1) == '+' {
				s.pos += 2
				s.token = ast.KindPlusPlusToken
			} else {
				s.pos++
				s.token = ast.KindPlusToken
			}
		case ',':
			s.pos++
			s.token = ast.KindCommaToken
		case '-':
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindMinusEqualsToken
			} else if s.charAt(1) == '-' {
				s.pos += 2
				s.token = ast.KindMinusMinusToken
			} else {
				s.pos++
				s.token = ast.KindMinusToken
			}
		case '.':
			if stringutil.IsDigit(s.charAt(1)) {
				s.token = s.scanNumber()
			} else if s.charAt(1) == '.' && s.charAt(2) == '.' {
				s.pos += 3
				s.token = ast.KindDotDotDotToken
			} else {
				s.pos++
				s.token = ast.KindDotToken
			}
		case '/':
			// Single-line comment
			if s.charAt(1) == '/' {
				s.pos += 2

				for {
					ch1, size := s.charAndSize()
					if size == 0 || stringutil.IsLineBreak(ch1) {
						break
					}
					s.pos += size
				}

				s.processCommentDirective(s.tokenStart, s.pos, false)

				if s.skipTrivia {
					continue
				}
				s.token = ast.KindSingleLineCommentTrivia
				return s.token
			}
			// Multi-line comment
			if s.charAt(1) == '*' {
				s.pos += 2
				isJSDoc := s.char() == '*' && s.charAt(1) != '/'

				commentClosed := false
				lastLineStart := s.tokenStart
				for {
					ch1, size := s.charAndSize()
					if size == 0 {
						break
					}

					if ch1 == '*' && s.charAt(1) == '/' {
						s.pos += 2
						commentClosed = true
						break
					}

					s.pos += size

					if stringutil.IsLineBreak(ch1) {
						lastLineStart = s.pos
						s.tokenFlags |= ast.TokenFlagsPrecedingLineBreak
					}
				}

				if isJSDoc && s.shouldParseJSDoc() {
					s.tokenFlags |= ast.TokenFlagsPrecedingJSDocComment
				}

				s.processCommentDirective(lastLineStart, s.pos, true)

				if !commentClosed {
					s.error(diagnostics.Asterisk_Slash_expected)
				}

				if s.skipTrivia {
					continue
				}

				if !commentClosed {
					s.tokenFlags |= ast.TokenFlagsUnterminated
				}
				s.token = ast.KindMultiLineCommentTrivia
				return s.token
			}
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindSlashEqualsToken
			} else {
				s.pos++
				s.token = ast.KindSlashToken
			}
		case '0':
			if s.charAt(1) == 'X' || s.charAt(1) == 'x' {
				start := s.pos
				s.pos += 2
				digits := s.scanHexDigits(1, true, true)
				if digits == "" {
					s.error(diagnostics.Hexadecimal_digit_expected)
					digits = "0"
				}
				if s.hexNumberCache == nil {
					s.hexNumberCache = make(map[string]string)
				}
				if cachedValue, ok := s.hexNumberCache[digits]; ok {
					s.tokenValue = cachedValue
				} else {
					rawText := s.text[start:s.pos]
					if strings.HasPrefix(rawText, "0x") && rawText[2:] == digits {
						s.tokenValue = rawText
					} else {
						s.tokenValue = "0x" + digits
					}
					s.hexNumberCache[digits] = s.tokenValue
				}
				s.tokenFlags |= ast.TokenFlagsHexSpecifier
				s.token = s.scanBigIntSuffix()
				break
			}
			if s.charAt(1) == 'B' || s.charAt(1) == 'b' {
				s.pos += 2
				digits := s.scanBinaryOrOctalDigits(2)
				if digits == "" {
					s.error(diagnostics.Binary_digit_expected)
					digits = "0"
				}
				s.tokenValue = "0b" + digits
				s.tokenFlags |= ast.TokenFlagsBinarySpecifier
				s.token = s.scanBigIntSuffix()
				break
			}
			if s.charAt(1) == 'O' || s.charAt(1) == 'o' {
				s.pos += 2
				digits := s.scanBinaryOrOctalDigits(8)
				if digits == "" {
					s.error(diagnostics.Octal_digit_expected)
					digits = "0"
				}
				s.tokenValue = "0o" + digits
				s.tokenFlags |= ast.TokenFlagsOctalSpecifier
				s.token = s.scanBigIntSuffix()
				break
			}
			fallthrough
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.token = s.scanNumber()
		case ':':
			s.pos++
			s.token = ast.KindColonToken
		case ';':
			s.pos++
			s.token = ast.KindSemicolonToken
		case '<':
			if isConflictMarkerTrivia(s.text, s.pos) {
				s.pos = scanConflictMarkerTrivia(s.text, s.pos, s.errorAt)
				if s.skipTrivia {
					continue
				} else {
					s.token = ast.KindConflictMarkerTrivia
					return s.token
				}
			}
			if s.charAt(1) == '<' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindLessThanLessThanEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindLessThanLessThanToken
				}
			} else if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindLessThanEqualsToken
			} else if s.languageVariant == core.LanguageVariantJSX && s.charAt(1) == '/' && s.charAt(2) != '*' {
				s.pos += 2
				s.token = ast.KindLessThanSlashToken
			} else {
				s.pos++
				s.token = ast.KindLessThanToken
			}
		case '=':
			if isConflictMarkerTrivia(s.text, s.pos) {
				s.pos = scanConflictMarkerTrivia(s.text, s.pos, s.errorAt)
				if s.skipTrivia {
					continue
				} else {
					s.token = ast.KindConflictMarkerTrivia
					return s.token
				}
			}
			if s.charAt(1) == '=' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindEqualsEqualsEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindEqualsEqualsToken
				}
			} else if s.charAt(1) == '>' {
				s.pos += 2
				s.token = ast.KindEqualsGreaterThanToken
			} else {
				s.pos++
				s.token = ast.KindEqualsToken
			}
		case '>':
			if isConflictMarkerTrivia(s.text, s.pos) {
				s.pos = scanConflictMarkerTrivia(s.text, s.pos, s.errorAt)
				if s.skipTrivia {
					continue
				} else {
					s.token = ast.KindConflictMarkerTrivia
					return s.token
				}
			}
			s.pos++
			s.token = ast.KindGreaterThanToken
		case '?':
			if s.charAt(1) == '.' && !stringutil.IsDigit(s.charAt(2)) {
				s.pos += 2
				s.token = ast.KindQuestionDotToken
			} else if s.charAt(1) == '?' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindQuestionQuestionEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindQuestionQuestionToken
				}
			} else {
				s.pos++
				s.token = ast.KindQuestionToken
			}
		case '[':
			s.pos++
			s.token = ast.KindOpenBracketToken
		case ']':
			s.pos++
			s.token = ast.KindCloseBracketToken
		case '^':
			if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindCaretEqualsToken
			} else {
				s.pos++
				s.token = ast.KindCaretToken
			}
		case '{':
			s.pos++
			s.token = ast.KindOpenBraceToken
		case '|':
			if isConflictMarkerTrivia(s.text, s.pos) {
				s.pos = scanConflictMarkerTrivia(s.text, s.pos, s.errorAt)
				if s.skipTrivia {
					continue
				} else {
					s.token = ast.KindConflictMarkerTrivia
					return s.token
				}
			}
			if s.charAt(1) == '|' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindBarBarEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindBarBarToken
				}
			} else if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindBarEqualsToken
			} else {
				s.pos++
				s.token = ast.KindBarToken
			}
		case '}':
			s.pos++
			s.token = ast.KindCloseBraceToken
		case '~':
			s.pos++
			s.token = ast.KindTildeToken
		case '@':
			s.pos++
			s.token = ast.KindAtToken
		case '\\':
			cp := s.peekUnicodeEscape()
			if cp >= 0 && IsIdentifierStart(cp) {
				s.tokenValue = string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
				s.token = GetIdentifierToken(s.tokenValue)
			} else {
				s.scanInvalidCharacter()
			}
		case '#':
			if s.charAt(1) == '!' {
				if s.pos == 0 {
					s.pos += 2
					for ch, size := s.charAndSize(); size > 0 && !stringutil.IsLineBreak(ch); ch, size = s.charAndSize() {
						s.pos += size
					}
					continue
				}
				s.errorAt(diagnostics.X_can_only_be_used_at_the_start_of_a_file, s.pos, 2)
				s.pos++
				s.token = ast.KindUnknown
				break
			}
			if s.charAt(1) == '\\' {
				s.pos++
				cp := s.peekUnicodeEscape()
				if cp >= 0 && IsIdentifierStart(cp) {
					s.tokenValue = "#" + string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
					s.token = ast.KindPrivateIdentifier
					break
				}
				s.pos--
			}
			if s.scanIdentifier(1) {
				s.token = ast.KindPrivateIdentifier
			} else {
				s.errorAt(diagnostics.Invalid_character, s.pos-1, 1)
				s.token = ast.KindUnknown
			}
		default:
			if ch < 0 {
				s.token = ast.KindEndOfFile
				break
			}
			if s.scanIdentifier(0) {
				s.token = GetIdentifierToken(s.tokenValue)
				break
			}
			ch, size := s.charAndSize()
			if ch == utf8.RuneError && size == 1 {
				s.errorAt(diagnostics.File_appears_to_be_binary, 0, 0)
				s.pos = len(s.text)
				s.token = ast.KindNonTextFileMarkerTrivia
				break
			}
			if stringutil.IsWhiteSpaceSingleLine(ch) {
				s.pos += size

				// If we get here and it's not 0x0085 (nextLine), then we're handling non-ASCII whitespace.
				// Handle skipTrivia like we do in the space case above.
				if ch == 0x0085 || s.skipTrivia {
					continue
				}

				for {
					ch, size = s.charAndSize()
					if !stringutil.IsWhiteSpaceSingleLine(ch) {
						break
					}
					s.pos += size
				}
				s.token = ast.KindWhitespaceTrivia
				return s.token
			}
			if stringutil.IsLineBreak(ch) {
				s.tokenFlags |= ast.TokenFlagsPrecedingLineBreak
				s.pos += size
				continue
			}
			s.scanInvalidCharacter()
		}
		return s.token
	}
}

func (s *Scanner) processCommentDirective(start int, end int, multiline bool) {
	// Skip starting slashes and whitespace
	pos := start
	if multiline {
		// Skip whitespace
		for pos < end && (s.text[pos] == ' ' || s.text[pos] == '\t') {
			pos++
		}
		// Skip combinations of / and *
		for pos < end && (s.text[pos] == '/' || s.text[pos] == '*') {
			pos++
		}
	} else {
		// Skip opening //
		pos += 2
		// Skip another / if present
		for pos < end && s.text[pos] == '/' {
			pos++
		}
	}
	// Skip whitespace
	for pos < end && (s.text[pos] == ' ' || s.text[pos] == '\t') {
		pos++
	}
	// Directive must start with '@'
	if !(pos < end && s.text[pos] == '@') {
		return
	}
	pos++
	var kind ast.CommentDirectiveKind
	switch {
	case strings.HasPrefix(s.text[pos:], "ts-expect-error"):
		kind = ast.CommentDirectiveKindExpectError
	case strings.HasPrefix(s.text[pos:], "ts-ignore"):
		kind = ast.CommentDirectiveKindIgnore
	default:
		return
	}
	s.commentDirectives = append(s.commentDirectives, ast.CommentDirective{Loc: core.NewTextRange(start, end), Kind: kind})
}

func (s *Scanner) ReScanLessThanToken() ast.Kind {
	if s.token == ast.KindLessThanLessThanToken {
		s.pos = s.tokenStart + 1
		s.token = ast.KindLessThanToken
	}
	return s.token
}

func (s *Scanner) ReScanGreaterThanToken() ast.Kind {
	if s.token == ast.KindGreaterThanToken {
		s.pos = s.tokenStart + 1
		if s.char() == '>' {
			if s.charAt(1) == '>' {
				if s.charAt(2) == '=' {
					s.pos += 3
					s.token = ast.KindGreaterThanGreaterThanGreaterThanEqualsToken
				} else {
					s.pos += 2
					s.token = ast.KindGreaterThanGreaterThanGreaterThanToken
				}
			} else if s.charAt(1) == '=' {
				s.pos += 2
				s.token = ast.KindGreaterThanGreaterThanEqualsToken
			} else {
				s.pos++
				s.token = ast.KindGreaterThanGreaterThanToken
			}
		} else if s.char() == '=' {
			s.pos++
			s.token = ast.KindGreaterThanEqualsToken
		}
	}
	return s.token
}

func (s *Scanner) ReScanTemplateToken(isTaggedTemplate bool) ast.Kind {
	s.pos = s.tokenStart
	s.token = s.scanTemplateAndSetTokenValue(!isTaggedTemplate)
	return s.token
}

func (s *Scanner) ReScanAsteriskEqualsToken() ast.Kind {
	if s.token != ast.KindAsteriskEqualsToken {
		panic("'ReScanAsteriskEqualsToken' should only be called on a '*='")
	}
	s.pos = s.tokenStart + 1
	s.token = ast.KindEqualsToken
	return s.token
}

// !!! https://github.com/microsoft/TypeScript/pull/55600
func (s *Scanner) ReScanSlashToken() ast.Kind {
	if s.token == ast.KindSlashToken || s.token == ast.KindSlashEqualsToken {
		s.pos = s.tokenStart + 1
		startOfRegExpBody := s.pos
		inEscape := false
		inCharacterClass := false
	loop:
		for {
			ch, size := s.charAndSize()
			// If we reach the end of a file, or hit a newline, then this is an unterminated
			// regex.  Report error and return what we have so far.
			switch {
			case size == 0 || stringutil.IsLineBreak(ch):
				s.tokenFlags |= ast.TokenFlagsUnterminated
				break loop
			case inEscape:
				// Parsing an escape character;
				// reset the flag and just advance to the next char.
				inEscape = false
			case ch == '/' && !inCharacterClass:
				// A slash within a character class is permissible,
				// but in general it signals the end of the regexp literal.
				break loop
			case ch == '[':
				inCharacterClass = true
			case ch == '\\':
				inEscape = true
			case ch == ']':
				inCharacterClass = false
			}
			s.pos += size
		}
		if s.tokenFlags&ast.TokenFlagsUnterminated != 0 {
			// Search for the nearest unbalanced bracket for better recovery. Since the expression is
			// invalid anyways, we take nested square brackets into consideration for the best guess.
			endOfRegExpBody := s.pos
			s.pos = startOfRegExpBody
			inEscape = false
			characterClassDepth := 0
			inDecimalQuantifier := false
			groupDepth := 0
			for s.pos < endOfRegExpBody {
				ch, size := s.charAndSize()
				if inEscape {
					inEscape = false
				} else if ch == '\\' {
					inEscape = true
				} else if ch == '[' {
					characterClassDepth++
				} else if ch == ']' && characterClassDepth != 0 {
					characterClassDepth--
				} else if characterClassDepth == 0 {
					if ch == '{' {
						inDecimalQuantifier = true
					} else if ch == '}' && inDecimalQuantifier {
						inDecimalQuantifier = false
					} else if !inDecimalQuantifier {
						if ch == '(' {
							groupDepth++
						} else if ch == ')' && groupDepth != 0 {
							groupDepth--
						} else if ch == ')' || ch == ']' || ch == '}' {
							// We encountered an unbalanced bracket outside a character class. Treat this position as the end of regex.
							break
						}
					}
				}
				s.pos += size
			}
			// Whitespaces and semicolons at the end are not likely to be part of the regex
			for {
				ch, size := utf8.DecodeLastRuneInString(s.text[:s.pos])
				if stringutil.IsWhiteSpaceLike(ch) || ch == ';' {
					s.pos -= size
				} else {
					break
				}
			}
			s.errorAt(diagnostics.Unterminated_regular_expression_literal, s.tokenStart, s.pos-s.tokenStart)
		} else {
			// Consume the slash character
			s.pos++
			for {
				ch, size := s.charAndSize()
				if size == 0 || !IsIdentifierPart(ch) {
					break
				}
				s.pos += size
			}
		}
		s.tokenValue = s.text[s.tokenStart:s.pos]
		s.token = ast.KindRegularExpressionLiteral
	}
	return s.token
}

func (s *Scanner) ReScanJsxToken(allowMultilineJsxText bool) ast.Kind {
	s.pos = s.fullStartPos
	s.tokenStart = s.fullStartPos
	s.token = s.ScanJsxTokenEx(allowMultilineJsxText)
	return s.token
}

func (s *Scanner) ReScanHashToken() ast.Kind {
	if s.token == ast.KindPrivateIdentifier {
		s.pos = s.tokenStart + 1
		s.token = ast.KindHashToken
	}
	return s.token
}

func (s *Scanner) ReScanQuestionToken() ast.Kind {
	if s.token != ast.KindQuestionQuestionToken {
		panic("'reScanQuestionToken' should only be called on a '??'")
	}
	s.pos = s.tokenStart + 1
	s.token = ast.KindQuestionToken
	return s.token
}

func (s *Scanner) ScanJsxToken() ast.Kind {
	return s.ScanJsxTokenEx(true /*allowMultilineJsxText*/)
}

func (s *Scanner) ScanJsxTokenEx(allowMultilineJsxText bool) ast.Kind {
	s.fullStartPos = s.pos
	s.tokenStart = s.pos
	ch := s.char()
	switch {
	case ch < 0:
		s.token = ast.KindEndOfFile
	case ch == '<':
		if s.charAt(1) == '/' {
			s.pos += 2
			s.token = ast.KindLessThanSlashToken
		} else {
			s.pos++
			s.token = ast.KindLessThanToken
		}
	case ch == '{':
		s.pos++
		s.token = ast.KindOpenBraceToken
	default:
		// First non-whitespace character on this line.
		firstNonWhitespace := 0
		// These initial values are special because the first line is:
		// firstNonWhitespace = 0 to indicate that we want leading whitespace
		for {
			ch, size := s.charAndSize()
			if size == 0 || ch == '{' {
				break
			}
			if ch == '<' {
				if isConflictMarkerTrivia(s.text, s.pos) {
					s.pos = scanConflictMarkerTrivia(s.text, s.pos, s.errorAt)
					s.token = ast.KindConflictMarkerTrivia
					return s.token
				}
				break
			}
			if ch == '>' {
				s.errorAt(diagnostics.Unexpected_token_Did_you_mean_or_gt, s.pos, 1)
			} else if ch == '}' {
				s.errorAt(diagnostics.Unexpected_token_Did_you_mean_or_rbrace, s.pos, 1)
			}
			// FirstNonWhitespace is 0, then we only see whitespaces so far. If we see a linebreak, we want to ignore that whitespaces.
			// i.e (- : whitespace)
			//      <div>----
			//      </div> becomes <div></div>
			//
			//      <div>----</div> becomes <div>----</div>
			if stringutil.IsLineBreak(ch) && firstNonWhitespace == 0 {
				firstNonWhitespace = -1
			} else if !allowMultilineJsxText && stringutil.IsLineBreak(ch) && firstNonWhitespace > 0 {
				// Stop JsxText on each line during formatting. This allows the formatter to
				// indent each line correctly.
				break
			} else if !stringutil.IsWhiteSpaceLike(ch) {
				firstNonWhitespace = s.pos
			}
			s.pos += size
		}
		s.tokenValue = s.text[s.fullStartPos:s.pos]
		s.token = ast.KindJsxText
		if firstNonWhitespace == -1 {
			s.token = ast.KindJsxTextAllWhiteSpaces
		}
	}
	return s.token
}

// Scans a JSX identifier; these differ from normal identifiers in that they allow dashes
func (s *Scanner) ScanJsxIdentifier() ast.Kind {
	if tokenIsIdentifierOrKeyword(s.token) {
		// An identifier or keyword has already been parsed - check for a `-` or a single instance of `:` and then append it and
		// everything after it to the token
		// Do note that this means that `scanJsxIdentifier` effectively _mutates_ the visible token without advancing to a new token
		// Any caller should be expecting this behavior and should only read the pos or token value after calling it.
		for {
			ch := s.char()
			if ch < 0 {
				break
			}
			if ch == '-' {
				s.tokenValue += "-"
				s.pos++
				continue
			}
			oldPos := s.pos
			s.tokenValue += s.scanIdentifierParts() // reuse `scanIdentifierParts` so unicode escapes are handled
			if s.pos == oldPos {
				break
			}
		}
		s.token = GetIdentifierToken(s.tokenValue)
	}
	return s.token
}

func (s *Scanner) ScanJsxAttributeValue() ast.Kind {
	s.fullStartPos = s.pos
	switch s.char() {
	case '"', '\'':
		s.tokenValue = s.scanString(true /*jsxAttributeString*/)
		s.token = ast.KindStringLiteral
		return s.token
	default:
		// If this scans anything other than `{`, it's a parse error.
		return s.Scan()
	}
}

func (s *Scanner) ReScanJsxAttributeValue() ast.Kind {
	s.pos = s.fullStartPos
	s.tokenStart = s.fullStartPos
	return s.ScanJsxAttributeValue()
}

/** In addition to the usual JSDoc ast.Kinds, can also return ast.KindJSDocCommentTextToken */
func (s *Scanner) ScanJSDocCommentTextToken(inBackticks bool) ast.Kind {
	s.fullStartPos = s.pos
	s.tokenFlags = ast.TokenFlagsNone
	if s.pos >= len(s.text) {
		s.token = ast.KindEndOfFile
		return s.token
	}
	s.tokenStart = s.pos
	for ch, size := s.charAndSize(); s.pos < len(s.text) && !stringutil.IsLineBreak(ch) && ch != '`'; ch, size = s.charAndSize() {
		if !inBackticks {
			if ch == '{' {
				break
			} else if ch == '@' && s.pos >= 0 {
				// @ doesn't start a new tag inside ``, and elsewhere, only after whitespace and before non-whitespace
				previous, _ := utf8.DecodeLastRuneInString(s.text[:s.pos])
				if stringutil.IsWhiteSpaceSingleLine(previous) {
					if s.pos+size >= len(s.text) {
						// EOF counts as non-whitespace
						break
					}
					next, _ := utf8.DecodeRuneInString(s.text[s.pos+size:])
					if !stringutil.IsWhiteSpaceLike(next) {
						break
					}
				}
			}
		}
		s.pos += size
	}
	if s.pos == s.tokenStart {
		return s.ScanJSDocToken()
	}
	s.tokenValue = s.text[s.tokenStart:s.pos]
	s.token = ast.KindJSDocCommentTextToken
	return s.token
}

func (s *Scanner) ScanJSDocToken() ast.Kind {
	s.fullStartPos = s.pos
	s.tokenFlags = ast.TokenFlagsNone
	if s.pos >= len(s.text) {
		s.token = ast.KindEndOfFile
		return s.token
	}

	s.tokenStart = s.pos
	ch, size := s.charAndSize()
	s.pos += size
	switch ch {
	case '\t', '\v', '\f', ' ':
		for ch2, size2 := s.charAndSize(); size2 > 0 && stringutil.IsWhiteSpaceSingleLine(ch2); ch2, size2 = s.charAndSize() {
			s.pos += size2
		}
		s.token = ast.KindWhitespaceTrivia
		return s.token
	case '@':
		s.token = ast.KindAtToken
		return s.token
	case '\r':
		if s.char() == '\n' {
			s.pos++
		}
		fallthrough
	case '\n':
		s.tokenFlags |= ast.TokenFlagsPrecedingLineBreak
		s.token = ast.KindNewLineTrivia
		return s.token
	case '*':
		s.token = ast.KindAsteriskToken
		return s.token
	case '{':
		s.token = ast.KindOpenBraceToken
		return s.token
	case '}':
		s.token = ast.KindCloseBraceToken
		return s.token
	case '[':
		s.token = ast.KindOpenBracketToken
		return s.token
	case ']':
		s.token = ast.KindCloseBracketToken
		return s.token
	case '(':
		s.token = ast.KindOpenParenToken
		return s.token
	case ')':
		s.token = ast.KindCloseParenToken
		return s.token
	case '<':
		s.token = ast.KindLessThanToken
		return s.token
	case '>':
		s.token = ast.KindGreaterThanToken
		return s.token
	case '=':
		s.token = ast.KindEqualsToken
		return s.token
	case ',':
		s.token = ast.KindCommaToken
		return s.token
	case '.':
		s.token = ast.KindDotToken
		return s.token
	case '`':
		s.token = ast.KindBacktickToken
		return s.token
	case '#':
		s.token = ast.KindHashToken
		return s.token
	case '\\':
		s.pos--
		cp := s.peekUnicodeEscape()
		if cp >= 0 && IsIdentifierStart(cp) {
			s.tokenValue = string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
			s.token = GetIdentifierToken(s.tokenValue)
		} else {
			s.scanInvalidCharacter()
		}
		return s.token
	}

	if IsIdentifierStart(ch) {
		char := ch
		for {
			if s.pos >= len(s.text) {
				break
			}
			char, size = s.charAndSize()
			if !IsIdentifierPart(char) && char != '-' {
				break
			}
			s.pos += size
		}
		s.tokenValue = s.text[s.tokenStart:s.pos]
		if char == '\\' {
			s.tokenValue += s.scanIdentifierParts()
		}
		s.token = GetIdentifierToken(s.tokenValue)
		return s.token
	} else {
		s.token = ast.KindUnknown
		return s.token
	}
}

func (s *Scanner) scanIdentifier(prefixLength int) bool {
	start := s.pos
	s.pos += prefixLength
	ch := s.char()
	// Fast path for simple ASCII identifiers
	if stringutil.IsASCIILetter(ch) || ch == '_' || ch == '$' {
		for {
			s.pos++
			ch = s.char()
			if !(isWordCharacter(ch) || ch == '$') {
				break
			}
		}
		if ch < utf8.RuneSelf && ch != '\\' {
			s.tokenValue = s.text[start:s.pos]
			return true
		}
		s.pos = start + prefixLength
	}
	ch, size := s.charAndSize()
	if IsIdentifierStart(ch) {
		for {
			s.pos += size
			ch, size = s.charAndSize()
			if !IsIdentifierPart(ch) {
				break
			}
		}
		s.tokenValue = s.text[start:s.pos]
		if ch == '\\' {
			s.tokenValue += s.scanIdentifierParts()
		}
		return true
	}
	return false
}

func (s *Scanner) scanIdentifierParts() string {
	var sb strings.Builder
	start := s.pos
	for {
		ch, size := s.charAndSize()
		if IsIdentifierPart(ch) {
			s.pos += size
			continue
		}
		if ch == '\\' {
			escaped := s.peekUnicodeEscape()
			if escaped >= 0 && IsIdentifierPart(escaped) {
				sb.WriteString(s.text[start:s.pos])
				sb.WriteRune(s.scanUnicodeEscape(true))
				start = s.pos
				continue
			}
		}
		break
	}
	sb.WriteString(s.text[start:s.pos])
	return sb.String()
}

func (s *Scanner) scanString(jsxAttributeString bool) string {
	quote := s.char()
	s.pos++
	// Fast path for simple strings without escape sequences.
	strLen := strings.IndexRune(s.text[s.pos:], quote)
	if strLen == 0 {
		s.pos++
		return ""
	}
	if strLen > 0 {
		str := s.text[s.pos : s.pos+strLen]
		if !jsxAttributeString && !strings.ContainsAny(str, "\r\n\\") {
			s.pos += strLen + 1
			return str
		}
	}
	var sb strings.Builder
	start := s.pos
	for {
		ch := s.char()
		if ch < 0 {
			sb.WriteString(s.text[start:s.pos])
			s.tokenFlags |= ast.TokenFlagsUnterminated
			s.error(diagnostics.Unterminated_string_literal)
			break
		}
		if ch == quote {
			sb.WriteString(s.text[start:s.pos])
			s.pos++
			break
		}
		if ch == '\\' && !jsxAttributeString {
			sb.WriteString(s.text[start:s.pos])
			sb.WriteString(s.scanEscapeSequence(EscapeSequenceScanningFlagsString | EscapeSequenceScanningFlagsReportErrors))
			start = s.pos
			continue
		}
		if (ch == '\n' || ch == '\r') && !jsxAttributeString {
			sb.WriteString(s.text[start:s.pos])
			s.tokenFlags |= ast.TokenFlagsUnterminated
			s.error(diagnostics.Unterminated_string_literal)
			break
		}
		s.pos++
	}
	return sb.String()
}

func (s *Scanner) scanTemplateAndSetTokenValue(shouldEmitInvalidEscapeError bool) ast.Kind {
	startedWithBacktick := s.char() == '`'
	s.pos++
	start := s.pos
	parts := make([]string, 0, 4)
	var token ast.Kind
	for {
		ch := s.char()
		if ch < 0 || ch == '`' {
			parts = append(parts, s.text[start:s.pos])
			if ch == '`' {
				s.pos++
			} else {
				s.tokenFlags |= ast.TokenFlagsUnterminated
				s.error(diagnostics.Unterminated_template_literal)
			}
			token = core.IfElse(startedWithBacktick, ast.KindNoSubstitutionTemplateLiteral, ast.KindTemplateTail)
			break
		}
		if ch == '$' && s.charAt(1) == '{' {
			parts = append(parts, s.text[start:s.pos])
			s.pos += 2
			token = core.IfElse(startedWithBacktick, ast.KindTemplateHead, ast.KindTemplateMiddle)
			break
		}
		if ch == '\\' {
			parts = append(parts, s.text[start:s.pos])
			parts = append(parts, s.scanEscapeSequence(EscapeSequenceScanningFlagsString|core.IfElse(shouldEmitInvalidEscapeError, EscapeSequenceScanningFlagsReportErrors, 0)))
			start = s.pos
			continue
		}
		// Speculated ECMAScript 6 Spec 11.8.6.1:
		// <CR><LF> and <CR> LineTerminatorSequences are normalized to <LF> for Template Values
		if ch == '\r' {
			parts = append(parts, s.text[start:s.pos])
			s.pos++
			if s.char() == '\n' {
				s.pos++
			}
			parts = append(parts, "\n")
			start = s.pos
			continue
		}
		s.pos++
	}
	s.tokenValue = strings.Join(parts, "")
	return token
}

func (s *Scanner) scanEscapeSequence(flags EscapeSequenceScanningFlags) string {
	start := s.pos
	s.pos++
	ch := s.char()
	if ch < 0 {
		s.error(diagnostics.Unexpected_end_of_text)
		return ""
	}
	s.pos++
	switch ch {
	case '0':
		// Although '0' preceding any digit is treated as LegacyOctalEscapeSequence,
		// '\08' should separately be interpreted as '\0' + '8'.
		if !stringutil.IsDigit(s.char()) {
			return "\x00"
		}
		// '\01', '\011'
		fallthrough
	case '1', '2', '3':
		// '\1', '\17', '\177'
		if stringutil.IsOctalDigit(s.char()) {
			s.pos++
		}
		// '\17', '\177'
		fallthrough
	case '4', '5', '6', '7':
		// '\4', '\47' but not '\477'
		if stringutil.IsOctalDigit(s.char()) {
			s.pos++
		}
		// '\47'
		s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
		if flags&EscapeSequenceScanningFlagsReportInvalidEscapeErrors != 0 {
			code, _ := strconv.ParseInt(s.text[start+1:s.pos], 8, 32)
			if flags&EscapeSequenceScanningFlagsRegularExpression != 0 && flags&EscapeSequenceScanningFlagsAtomEscape == 0 && ch != '0' {
				s.errorAt(diagnostics.Octal_escape_sequences_and_backreferences_are_not_allowed_in_a_character_class_If_this_was_intended_as_an_escape_sequence_use_the_syntax_0_instead, start, s.pos-start, fmt.Sprintf("%02x", code))
			} else {
				s.errorAt(diagnostics.Octal_escape_sequences_are_not_allowed_Use_the_syntax_0, start, s.pos-start, "\\x"+fmt.Sprintf("%02x", code))
			}
			return string(rune(code))
		}
		return s.text[start:s.pos]
	case '8', '9':
		// the invalid '\8' and '\9'
		s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
		if flags&EscapeSequenceScanningFlagsReportInvalidEscapeErrors != 0 {
			if flags&EscapeSequenceScanningFlagsRegularExpression != 0 && flags&EscapeSequenceScanningFlagsAtomEscape == 0 {
				s.errorAt(diagnostics.Decimal_escape_sequences_and_backreferences_are_not_allowed_in_a_character_class, start, s.pos-start)
			} else {
				s.errorAt(diagnostics.Escape_sequence_0_is_not_allowed, start, s.pos-start, s.text[start:s.pos])
			}
			return string(ch)
		}
		return s.text[start:s.pos]
	case 'b':
		return "\b"
	case 't':
		return "\t"
	case 'n':
		return "\n"
	case 'v':
		return "\v"
	case 'f':
		return "\f"
	case 'r':
		return "\r"
	case '\'':
		return "'"
	case '"':
		return "\""
	case 'u':
		// '\uDDDD' and '\U{DDDDDD}'
		s.pos -= 2
		codePoint := s.scanUnicodeEscape(flags&EscapeSequenceScanningFlagsReportInvalidEscapeErrors != 0)
		if codePoint < 0 {
			return s.text[start:s.pos]
		}
		return string(codePoint)
	case 'x':
		// '\xDD'
		for ; s.pos < start+4; s.pos++ {
			if !stringutil.IsHexDigit(s.char()) {
				s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
				if flags&EscapeSequenceScanningFlagsReportInvalidEscapeErrors != 0 {
					s.error(diagnostics.Hexadecimal_digit_expected)
				}
				return s.text[start:s.pos]
			}
		}
		s.tokenFlags |= ast.TokenFlagsHexEscape
		escapedValue, _ := strconv.ParseInt(s.text[start+2:s.pos], 16, 32)
		return string(rune(escapedValue))
	case '\r':
		// when encountering a LineContinuation (i.e. a backslash and a line terminator sequence),
		// the line terminator is interpreted to be "the empty code unit sequence".
		if s.char() == '\n' {
			s.pos++
		}
		fallthrough
	case '\n':
		// case CharacterCodes.lineSeparator !!!
		// case CharacterCodes.paragraphSeparator !!!
		return ""
	default:
		if flags&EscapeSequenceScanningFlagsAnyUnicodeMode != 0 || flags&EscapeSequenceScanningFlagsRegularExpression != 0 && flags&EscapeSequenceScanningFlagsAnnexB == 0 && IsIdentifierPart(ch) {
			s.errorAt(diagnostics.This_character_cannot_be_escaped_in_a_regular_expression, s.pos-2, 2)
		}
		return string(ch)
	}
}

// Known to be at \u
func (s *Scanner) scanUnicodeEscape(shouldEmitInvalidEscapeError bool) rune {
	s.pos += 2
	start := s.pos
	extended := s.char() == '{'
	var hexDigits string
	if extended {
		s.pos++
		s.tokenFlags |= ast.TokenFlagsExtendedUnicodeEscape
		hexDigits = s.scanHexDigits(1, true, false)
	} else {
		s.tokenFlags |= ast.TokenFlagsUnicodeEscape
		hexDigits = s.scanHexDigits(4, false, false)
	}
	if hexDigits == "" {
		s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
		if shouldEmitInvalidEscapeError {
			s.error(diagnostics.Hexadecimal_digit_expected)
		}
		return -1
	}
	hexValue, _ := strconv.ParseInt(hexDigits, 16, 32)
	if extended {
		if hexValue > 0x10FFFF {
			s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
			if shouldEmitInvalidEscapeError {
				s.errorAt(diagnostics.An_extended_Unicode_escape_value_must_be_between_0x0_and_0x10FFFF_inclusive, start+1, s.pos-start-1)
			}
			return -1
		}
		if s.char() != '}' {
			s.tokenFlags |= ast.TokenFlagsContainsInvalidEscape
			if shouldEmitInvalidEscapeError {
				s.error(diagnostics.Unterminated_Unicode_escape_sequence)
			}
			return -1
		}
		s.pos++
	}
	return rune(hexValue)
}

// Current character is known to be a backslash. Check for Unicode escape of the form '\uXXXX'
// or '\u{XXXXXX}' and return code point value if valid Unicode escape is found. Otherwise return -1.
func (s *Scanner) peekUnicodeEscape() rune {
	if s.charAt(1) == 'u' {
		savePos := s.pos
		saveTokenFlags := s.tokenFlags
		codePoint := s.scanUnicodeEscape(false)
		s.pos = savePos
		s.tokenFlags = saveTokenFlags
		return codePoint
	}
	return -1
}

func (s *Scanner) scanNumber() ast.Kind {
	start := s.pos
	var fixedPart string
	if s.char() == '0' {
		s.pos++
		if s.char() == '_' {
			s.tokenFlags |= ast.TokenFlagsContainsSeparator | ast.TokenFlagsContainsInvalidSeparator
			s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos, 1)
			s.pos = start
			fixedPart = s.scanNumberFragment()
		} else {
			digits, isOctal := s.scanDigits()
			if digits == "" {
				fixedPart = "0"
			} else if !isOctal {
				s.tokenFlags |= ast.TokenFlagsContainsLeadingZero
				fixedPart = digits
			} else {
				s.tokenValue = jsnum.FromString(digits).String()
				s.tokenFlags |= ast.TokenFlagsOctal
				withMinus := s.token == ast.KindMinusToken
				literal := core.IfElse(withMinus, "-", "") + "0o" + s.tokenValue
				if withMinus {
					start--
				}
				s.errorAt(diagnostics.Octal_literals_are_not_allowed_Use_the_syntax_0, start, s.pos-start, literal)
				return ast.KindNumericLiteral
			}
		}
	} else {
		fixedPart = s.scanNumberFragment()
	}
	fixedPartEnd := s.pos
	fractionalPart := ""
	exponentPreamble := ""
	exponentPart := ""
	if s.char() == '.' {
		s.pos++
		fractionalPart = s.scanNumberFragment()
	}
	end := s.pos
	if s.char() == 'E' || s.char() == 'e' {
		s.pos++
		s.tokenFlags |= ast.TokenFlagsScientific
		if s.char() == '+' || s.char() == '-' {
			s.pos++
		}
		startNumericPart := s.pos
		exponentPart = s.scanNumberFragment()
		if exponentPart == "" {
			s.error(diagnostics.Digit_expected)
		} else {
			exponentPreamble = s.text[end:startNumericPart]
			end = s.pos
		}
	}
	if s.tokenFlags&ast.TokenFlagsContainsSeparator != 0 {
		s.tokenValue = fixedPart
		if fractionalPart != "" {
			s.tokenValue += "." + fractionalPart
		}
		if exponentPart != "" {
			s.tokenValue += exponentPreamble + exponentPart
		}
	} else {
		s.tokenValue = s.text[start:end]
	}
	if s.tokenFlags&ast.TokenFlagsContainsLeadingZero != 0 {
		s.errorAt(diagnostics.Decimals_with_leading_zeros_are_not_allowed, start, s.pos-start)
		s.tokenValue = jsnum.FromString(s.tokenValue).String()
		return ast.KindNumericLiteral
	}
	var result ast.Kind
	if fixedPartEnd == s.pos {
		result = s.scanBigIntSuffix()
	} else {
		s.tokenValue = jsnum.FromString(s.tokenValue).String()
		result = ast.KindNumericLiteral
	}
	ch, _ := s.charAndSize()
	if IsIdentifierStart(ch) {
		idStart := s.pos
		id := s.scanIdentifierParts()
		if result != ast.KindBigIntLiteral && len(id) == 1 && s.text[idStart] == 'n' {
			if s.tokenFlags&ast.TokenFlagsScientific != 0 {
				s.errorAt(diagnostics.A_bigint_literal_cannot_use_exponential_notation, start, s.pos-start)
				return result
			}
			if fixedPartEnd < idStart {
				s.errorAt(diagnostics.A_bigint_literal_must_be_an_integer, start, s.pos-start)
				return result
			}
		}
		s.errorAt(diagnostics.An_identifier_or_keyword_cannot_immediately_follow_a_numeric_literal, idStart, s.pos-idStart)
		s.pos = idStart
	}
	return result
}

func (s *Scanner) scanNumberFragment() string {
	start := s.pos
	allowSeparator := false
	isPreviousTokenSeparator := false
	result := ""
	for {
		ch := s.char()
		if ch == '_' {
			s.tokenFlags |= ast.TokenFlagsContainsSeparator
			if allowSeparator {
				allowSeparator = false
				isPreviousTokenSeparator = true
				result += s.text[start:s.pos]
			} else {
				s.tokenFlags |= ast.TokenFlagsContainsInvalidSeparator
				if isPreviousTokenSeparator {
					s.errorAt(diagnostics.Multiple_consecutive_numeric_separators_are_not_permitted, s.pos, 1)
				} else {
					s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos, 1)
				}
			}
			s.pos++
			start = s.pos
			continue
		}
		if stringutil.IsDigit(ch) {
			allowSeparator = true
			isPreviousTokenSeparator = false
			s.pos++
			continue
		}
		break
	}
	if isPreviousTokenSeparator {
		s.tokenFlags |= ast.TokenFlagsContainsInvalidSeparator
		s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos-1, 1)
	}
	return result + s.text[start:s.pos]
}

func (s *Scanner) scanDigits() (string, bool) {
	start := s.pos
	isOctal := true
	for stringutil.IsDigit(s.char()) {
		if !stringutil.IsOctalDigit(s.char()) {
			isOctal = false
		}
		s.pos++
	}
	return s.text[start:s.pos], isOctal
}

func (s *Scanner) scanHexDigits(minCount int, scanAsManyAsPossible bool, canHaveSeparators bool) string {
	digitCount := 0
	start := s.pos
	allowSeparator := false
	isPreviousTokenSeparator := false
	for digitCount < minCount || scanAsManyAsPossible {
		ch := s.char()
		if stringutil.IsHexDigit(ch) {
			allowSeparator = canHaveSeparators
			isPreviousTokenSeparator = false
			digitCount++
		} else if canHaveSeparators && ch == '_' {
			s.tokenFlags |= ast.TokenFlagsContainsSeparator
			if allowSeparator {
				allowSeparator = false
				isPreviousTokenSeparator = true
			} else if isPreviousTokenSeparator {
				s.errorAt(diagnostics.Multiple_consecutive_numeric_separators_are_not_permitted, s.pos, 1)
			} else {
				s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos, 1)
			}
		} else {
			break
		}
		s.pos++
	}
	if isPreviousTokenSeparator {
		s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos-1, 1)
	}
	if digitCount < minCount {
		return ""
	}
	digits := s.text[start:s.pos]
	if s.hexDigitCache == nil {
		s.hexDigitCache = make(map[string]string)
	}
	if cached, ok := s.hexDigitCache[digits]; ok {
		return cached
	} else {
		original := digits
		if s.tokenFlags&ast.TokenFlagsContainsSeparator != 0 {
			digits = strings.ReplaceAll(digits, "_", "")
		}
		digits = strings.ToLower(digits) // standardize hex literals to lowercase
		s.hexDigitCache[original] = digits
		return digits
	}
}

func (s *Scanner) scanBinaryOrOctalDigits(base int32) string {
	var sb strings.Builder
	allowSeparator := false
	isPreviousTokenSeparator := false
	for {
		ch := s.char()
		if stringutil.IsDigit(ch) && ch-'0' < base {
			sb.WriteByte(byte(ch))
			allowSeparator = true
			isPreviousTokenSeparator = false
		} else if ch == '_' {
			s.tokenFlags |= ast.TokenFlagsContainsSeparator
			if allowSeparator {
				allowSeparator = false
				isPreviousTokenSeparator = true
			} else if isPreviousTokenSeparator {
				s.errorAt(diagnostics.Multiple_consecutive_numeric_separators_are_not_permitted, s.pos, 1)
			} else {
				s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos, 1)
			}
		} else {
			break
		}
		s.pos++
	}
	if isPreviousTokenSeparator {
		s.errorAt(diagnostics.Numeric_separators_are_not_allowed_here, s.pos-1, 1)
	}
	return sb.String()
}

func (s *Scanner) scanBigIntSuffix() ast.Kind {
	if s.char() == 'n' {
		s.tokenValue += "n"
		if s.tokenFlags&ast.TokenFlagsBinaryOrOctalSpecifier != 0 {
			s.tokenValue = jsnum.ParsePseudoBigInt(s.tokenValue) + "n"
		}
		s.pos++
		return ast.KindBigIntLiteral
	}
	if s.numberCache == nil {
		s.numberCache = make(map[string]string)
	}
	if cached, ok := s.numberCache[s.tokenValue]; ok {
		s.tokenValue = cached
	} else {
		tokenValue := jsnum.FromString(s.tokenValue).String()
		if tokenValue == s.tokenValue {
			tokenValue = s.tokenValue
		}
		s.numberCache[s.tokenValue] = tokenValue
		s.tokenValue = tokenValue
	}
	return ast.KindNumericLiteral
}

func (s *Scanner) scanInvalidCharacter() {
	_, size := s.charAndSize()
	s.errorAt(diagnostics.Invalid_character, s.pos, size)
	s.pos += size
	s.token = ast.KindUnknown
}

func GetIdentifierToken(str string) ast.Kind {
	if len(str) >= 2 && len(str) <= 12 && str[0] >= 'a' && str[0] <= 'z' {
		keyword := textToKeyword[str]
		if keyword != ast.KindUnknown {
			return keyword
		}
	}
	return ast.KindIdentifier
}

func IsValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, ch := range s {
		if i == 0 && !IsIdentifierStart(ch) || i != 0 && !IsIdentifierPart(ch) {
			return false
		}
	}
	return true
}

// Section 6.1.4
func isWordCharacter(ch rune) bool {
	return stringutil.IsASCIILetter(ch) || stringutil.IsDigit(ch) || ch == '_'
}

func IsIdentifierStart(ch rune) bool {
	return stringutil.IsASCIILetter(ch) || ch == '_' || ch == '$' || ch >= utf8.RuneSelf && isUnicodeIdentifierStart(ch)
}

func IsIdentifierPart(ch rune) bool {
	return IsIdentifierPartEx(ch, core.LanguageVariantStandard)
}

func IsIdentifierPartEx(ch rune, languageVariant core.LanguageVariant) bool {
	return isWordCharacter(ch) || ch == '$' ||
		ch >= utf8.RuneSelf && isUnicodeIdentifierPart(ch) ||
		languageVariant == core.LanguageVariantJSX && (ch == '-' || ch == ':') // "-" and ":" are valid in JSX Identifiers
}

func isUnicodeIdentifierStart(ch rune) bool {
	return isInUnicodeRanges(ch, unicodeESNextIdentifierStart)
}

func isUnicodeIdentifierPart(ch rune) bool {
	return isInUnicodeRanges(ch, unicodeESNextIdentifierPart)
}

func isInUnicodeRanges(cp rune, ranges []rune) bool {
	// Bail out quickly if it couldn't possibly be in the map
	if cp < ranges[0] {
		return false
	}
	// Perform binary search in one of the Unicode range maps
	lo := 0
	hi := len(ranges)
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		// mid has to be even to catch beginning of a range
		mid -= mid % 2
		if ranges[mid] <= cp && cp <= ranges[mid+1] {
			return true
		}
		if cp < ranges[mid] {
			hi = mid
		} else {
			lo = mid + 2
		}
	}
	return false
}

var tokenToText map[ast.Kind]string

func init() {
	tokenToText = make(map[ast.Kind]string, len(textToToken))
	for text, key := range textToToken {
		tokenToText[key] = text
	}
}

func TokenToString(token ast.Kind) string {
	return tokenToText[token]
}

func StringToToken(s string) ast.Kind {
	kind, ok := textToToken[s]
	if ok {
		return kind
	}
	return ast.KindUnknown
}

func GetViableKeywordSuggestions() []string {
	result := make([]string, 0, len(textToKeyword))
	for text := range textToKeyword {
		if len(text) > 2 {
			result = append(result, text)
		}
	}
	return result
}

func couldStartTrivia(text string, pos int) bool {
	// Keep in sync with skipTrivia
	switch ch := text[pos]; ch {
	// Characters that could start normal trivia
	case '\r', '\n', '\t', '\v', '\f', ' ', '/',
		// Characters that could start conflict marker trivia
		'<', '|', '=', '>':
		return true
	case '#':
		// Only if its the beginning can we have #! trivia
		return pos == 0
	default:
		return ch > maxAsciiCharacter
	}
}

type SkipTriviaOptions struct {
	StopAfterLineBreak bool
	StopAtComments     bool
	InJSDoc            bool
}

func SkipTrivia(text string, pos int) int {
	return SkipTriviaEx(text, pos, nil)
}

func SkipTriviaEx(text string, pos int, options *SkipTriviaOptions) int {
	if ast.PositionIsSynthesized(pos) {
		return pos
	}
	if options == nil {
		options = &SkipTriviaOptions{}
	}

	textLen := len(text)
	canConsumeStar := false
	// Keep in sync with couldStartTrivia
	for {
		ch, size := utf8.DecodeRuneInString(text[pos:])
		switch ch {
		case '\r':
			if pos+1 < textLen && text[pos+1] == '\n' {
				pos++
			}
			fallthrough
		case '\n':
			pos++
			if options.StopAfterLineBreak {
				return pos
			}
			canConsumeStar = options.InJSDoc
			continue
		case '\t', '\v', '\f', ' ':
			pos++
			continue
		case '/':
			if options.StopAtComments {
				break
			}
			if pos+1 < textLen {
				if text[pos+1] == '/' {
					pos += 2
					for pos < textLen {
						ch, size := utf8.DecodeRuneInString(text[pos:])
						if stringutil.IsLineBreak(ch) {
							break
						}
						pos += size
					}
					canConsumeStar = false
					continue
				}
				if text[pos+1] == '*' {
					pos += 2
					for pos < textLen {
						if text[pos] == '*' && (pos+1 < textLen) && text[pos+1] == '/' {
							pos += 2
							break
						}
						_, size := utf8.DecodeRuneInString(text[pos:])
						pos += size
					}
					canConsumeStar = false
					continue
				}
			}
		case '<', '|', '=', '>':
			if isConflictMarkerTrivia(text, pos) {
				pos = scanConflictMarkerTrivia(text, pos, nil)
				canConsumeStar = false
				continue
			}
		case '#':
			if pos == 0 && isShebangTrivia(text, pos) {
				pos = scanShebangTrivia(text, pos)
				canConsumeStar = false
				continue
			}
		case '*':
			if canConsumeStar {
				pos++
				canConsumeStar = false
				continue
			}
		default:
			if ch > rune(maxAsciiCharacter) && stringutil.IsWhiteSpaceLike(ch) {
				pos += size
				continue
			}
		}
		return pos
	}
}

// All conflict markers consist of the same character repeated seven times.  If it is
// a <<<<<<< or >>>>>>> marker then it is also followed by a space.
var (
	mergeConflictMarkerLength      = len("<<<<<<<")
	maxAsciiCharacter         byte = 127
)

func isConflictMarkerTrivia(text string, pos int) bool {
	if pos < 0 {
		panic("pos < 0")
	}

	// Conflict markers must be at the start of a line.
	var prev rune
	if pos >= 2 {
		prev, _ = utf8.DecodeLastRuneInString(text[:pos-2])
	}
	if pos == 0 || stringutil.IsLineBreak(prev) || pos >= 1 && stringutil.IsLineBreak(rune(text[pos-1])) {
		ch := text[pos]

		if (pos + mergeConflictMarkerLength) < len(text) {
			for i := range mergeConflictMarkerLength {
				if text[pos+i] != ch {
					return false
				}
			}

			return ch == '=' || text[pos+mergeConflictMarkerLength] == ' '
		}
	}

	return false
}

func scanConflictMarkerTrivia(text string, pos int, reportError func(diag *diagnostics.Message, pos int, length int, args ...any)) int {
	if reportError != nil {
		reportError(diagnostics.Merge_conflict_marker_encountered, pos, mergeConflictMarkerLength)
	}
	ch, size := utf8.DecodeRuneInString(text[pos:])
	length := len(text)

	if ch == '<' || ch == '>' {
		for pos < length && !stringutil.IsLineBreak(ch) {
			pos += size
			ch, size = utf8.DecodeRuneInString(text[pos:])
		}
	} else {
		if ch != '|' && ch != '=' {
			panic("Assertion failed: ch must be either '|' or '='")
		}
		// Consume everything from the start of a ||||||| or ======= marker to the start
		// of the next ======= or >>>>>>> marker.
		for pos < length {
			currentChar := text[pos]
			if (currentChar == '=' || currentChar == '>') && rune(currentChar) != ch && isConflictMarkerTrivia(text, pos) {
				break
			}

			pos++
		}
	}

	return pos
}

func isShebangTrivia(text string, pos int) bool {
	if len(text) < 2 {
		return false
	}
	if pos != 0 {
		panic("Shebangs check must only be done at the start of the file")
	}
	return text[0] == '#' && text[1] == '!'
}

func scanShebangTrivia(text string, pos int) int {
	pos += 2
	for pos < len(text) {
		ch, size := utf8.DecodeRuneInString(text[pos:])
		if stringutil.IsLineBreak(ch) {
			break
		}
		pos += size
	}
	return pos
}

func GetShebang(text string) string {
	if !isShebangTrivia(text, 0) {
		return ""
	}

	end := scanShebangTrivia(text, 0)
	return text[:end]
}

func GetScannerForSourceFile(sourceFile *ast.SourceFile, pos int) *Scanner {
	s := NewScanner()
	s.text = sourceFile.Text()
	s.pos = pos
	s.languageVariant = sourceFile.LanguageVariant
	s.Scan()
	return s
}

func ScanTokenAtPosition(sourceFile *ast.SourceFile, pos int) ast.Kind {
	s := GetScannerForSourceFile(sourceFile, pos)
	return s.token
}

func GetRangeOfTokenAtPosition(sourceFile *ast.SourceFile, pos int) core.TextRange {
	s := GetScannerForSourceFile(sourceFile, pos)
	return core.NewTextRange(s.tokenStart, s.pos)
}

func GetTokenPosOfNode(node *ast.Node, sourceFile *ast.SourceFile, includeJSDoc bool) int {
	// With nodes that have no width (i.e. 'Missing' nodes), we actually *don't*
	// want to skip trivia because this will launch us forward to the next token.
	if ast.NodeIsMissing(node) {
		return node.Pos()
	}
	if ast.IsJSDocNode(node) || node.Kind == ast.KindJsxText {
		// JsxText cannot actually contain comments, even though the scanner will think it sees comments
		return SkipTriviaEx(sourceFile.Text(), node.Pos(), &SkipTriviaOptions{StopAtComments: true})
	}
	if includeJSDoc && len(node.JSDoc(sourceFile)) > 0 {
		return GetTokenPosOfNode(node.JSDoc(sourceFile)[0], sourceFile, false /*includeJSDoc*/)
	}
	return SkipTriviaEx(sourceFile.Text(), node.Pos(), &SkipTriviaOptions{InJSDoc: node.Flags&ast.NodeFlagsJSDoc != 0})
}

func getErrorRangeForArrowFunction(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	pos := SkipTrivia(sourceFile.Text(), node.Pos())
	body := node.AsArrowFunction().Body
	if body != nil && body.Kind == ast.KindBlock {
		startLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.Pos())
		endLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.End())
		if startLine < endLine {
			// The arrow function spans multiple lines,
			// make the error span be the first line, inclusive.
			return core.NewTextRange(pos, GetEndLinePosition(sourceFile, startLine))
		}
	}
	return core.NewTextRange(pos, node.End())
}

func GetErrorRangeForNode(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	errorNode := node
	switch node.Kind {
	case ast.KindSourceFile:
		pos := SkipTrivia(sourceFile.Text(), 0)
		if pos == len(sourceFile.Text()) {
			return core.NewTextRange(0, 0)
		}
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	// This list is a work in progress. Add missing node kinds to improve their error spans
	case ast.KindFunctionDeclaration, ast.KindMethodDeclaration:
		if node.Flags&ast.NodeFlagsReparsed != 0 {
			errorNode = node
			break
		}
		fallthrough
	case ast.KindVariableDeclaration, ast.KindBindingElement, ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration,
		ast.KindModuleDeclaration, ast.KindEnumDeclaration, ast.KindEnumMember, ast.KindFunctionExpression,
		ast.KindGetAccessor, ast.KindSetAccessor, ast.KindTypeAliasDeclaration, ast.KindJSTypeAliasDeclaration, ast.KindPropertyDeclaration,
		ast.KindPropertySignature, ast.KindNamespaceImport:
		errorNode = ast.GetNameOfDeclaration(node)
	case ast.KindArrowFunction:
		return getErrorRangeForArrowFunction(sourceFile, node)
	case ast.KindCaseClause, ast.KindDefaultClause:
		start := SkipTrivia(sourceFile.Text(), node.Pos())
		end := node.End()
		statements := node.AsCaseOrDefaultClause().Statements.Nodes
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return core.NewTextRange(start, end)
	case ast.KindReturnStatement, ast.KindYieldExpression:
		pos := SkipTrivia(sourceFile.Text(), node.Pos())
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindSatisfiesExpression:
		pos := SkipTrivia(sourceFile.Text(), node.AsSatisfiesExpression().Expression.End())
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindConstructor:
		if node.Flags&ast.NodeFlagsReparsed != 0 {
			errorNode = node
			break
		}
		scanner := GetScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.TokenStart()
		for scanner.Token() != ast.KindConstructorKeyword && scanner.Token() != ast.KindStringLiteral && scanner.Token() != ast.KindEndOfFile {
			scanner.Scan()
		}
		return core.NewTextRange(start, scanner.TokenEnd())
		// !!!
		// case KindJSDocSatisfiesTag:
		// 	pos := scanner.SkipTrivia(sourceFile.Text(), node.tagName.pos)
		// 	return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	}
	if errorNode == nil {
		// If we don't have a better node, then just set the error on the first token of
		// construct.
		return GetRangeOfTokenAtPosition(sourceFile, node.Pos())
	}
	pos := errorNode.Pos()
	if !ast.NodeIsMissing(errorNode) && !ast.IsJsxText(errorNode) {
		pos = SkipTrivia(sourceFile.Text(), pos)
	}
	return core.NewTextRange(pos, errorNode.End())
}

func ComputeLineOfPosition(lineStarts []core.TextPos, pos int) int {
	low := 0
	high := len(lineStarts) - 1
	for low <= high {
		middle := low + ((high - low) >> 1)
		value := int(lineStarts[middle])
		if value < pos {
			low = middle + 1
		} else if value > pos {
			high = middle - 1
		} else {
			return middle
		}
	}
	return low - 1
}

func GetLineStarts(sourceFile ast.SourceFileLike) []core.TextPos {
	return sourceFile.LineMap()
}

func GetLineAndCharacterOfPosition(sourceFile ast.SourceFileLike, pos int) (line int, character int) {
	lineMap := GetLineStarts(sourceFile)
	line = ComputeLineOfPosition(lineMap, pos)
	character = utf8.RuneCountInString(sourceFile.Text()[lineMap[line]:pos])
	return
}

func GetEndLinePosition(sourceFile *ast.SourceFile, line int) int {
	pos := int(GetLineStarts(sourceFile)[line])
	for {
		ch, size := utf8.DecodeRuneInString(sourceFile.Text()[pos:])
		if size == 0 || stringutil.IsLineBreak(ch) {
			return pos
		}
		pos += size
	}
}

func GetPositionOfLineAndCharacter(sourceFile *ast.SourceFile, line int, character int) int {
	return ComputePositionOfLineAndCharacter(GetLineStarts(sourceFile), line, character)
}

func ComputePositionOfLineAndCharacter(lineStarts []core.TextPos, line int, character int) int {
	/// !!! debugText, allowEdits
	if line < 0 || line >= len(lineStarts) {
		// if (allowEdits) {
		//     // Clamp line to nearest allowable value
		//     line = line < 0 ? 0 : line >= lineStarts.length ? lineStarts.length - 1 : line;
		// }
		panic(fmt.Sprintf("Bad line number. Line: %d, lineStarts.length: %d.", line, len(lineStarts)))
	}

	res := int(lineStarts[line]) + character

	// !!!
	// if (allowEdits) {
	//     // Clamp to nearest allowable values to allow the underlying to be edited without crashing (accuracy is lost, instead)
	//     // TODO: Somehow track edits between file as it was during the creation of sourcemap we have and the current file and
	//     // apply them to the computed position to improve accuracy
	//     return res > lineStarts[line + 1] ? lineStarts[line + 1] : typeof debugText === "string" && res > debugText.length ? debugText.length : res;
	// }
	if line < len(lineStarts)-1 && res >= int(lineStarts[line+1]) {
		panic("Computed position is beyond that of the following line.")
	}
	// !!!
	// else if (debugText !== undefined) {
	//     Debug.assert(res <= debugText.length); // Allow single character overflow for trailing newline
	// }
	return res
}

func GetLeadingCommentRanges(f *ast.NodeFactory, text string, pos int) iter.Seq[ast.CommentRange] {
	return iterateCommentRanges(f, text, pos, false)
}

func GetTrailingCommentRanges(f *ast.NodeFactory, text string, pos int) iter.Seq[ast.CommentRange] {
	return iterateCommentRanges(f, text, pos, true)
}

/*
Returns an iterator over each comment range following the provided position.
Single-line comment ranges include the leading double-slash characters but not the ending
line break. Multi-line comment ranges include the leading slash-asterisk and trailing
asterisk-slash characters.
*/
func iterateCommentRanges(f *ast.NodeFactory, text string, pos int, trailing bool) iter.Seq[ast.CommentRange] {
	return func(yield func(ast.CommentRange) bool) {
		var pendingPos int
		var pendingEnd int
		var pendingKind ast.Kind
		var pendingHasTrailingNewLine bool
		hasPendingCommentRange := false
		collecting := trailing
		if pos == 0 {
			collecting = true
			if isShebangTrivia(text, pos) {
				pos = scanShebangTrivia(text, pos)
			}
		}
	scan:
		for pos >= 0 && pos < len(text) {
			ch, size := utf8.DecodeRuneInString(text[pos:])
			switch ch {
			case '\r':
				if pos+1 < len(text) && text[pos+1] == '\n' {
					pos++
				}
				fallthrough
			case '\n':
				pos++
				if trailing {
					break scan
				}

				collecting = true
				if hasPendingCommentRange {
					pendingHasTrailingNewLine = true
				}

				continue
			case '\t', '\v', '\f', ' ':
				pos++
				continue
			case '/':
				var nextChar byte
				if pos+1 < len(text) {
					nextChar = text[pos+1]
				}
				hasTrailingNewLine := false
				if nextChar == '/' || nextChar == '*' {
					var kind ast.Kind
					if nextChar == '/' {
						kind = ast.KindSingleLineCommentTrivia
					} else {
						kind = ast.KindMultiLineCommentTrivia
					}

					startPos := pos
					pos += 2
					if nextChar == '/' {
						for pos < len(text) {
							c, s := utf8.DecodeRuneInString(text[pos:])
							if stringutil.IsLineBreak(c) {
								hasTrailingNewLine = true
								break
							}
							pos += s
						}
					} else {
						for pos < len(text) {
							c, s := utf8.DecodeRuneInString(text[pos:])
							if c == '*' && pos+1 < len(text) && text[pos+1] == '/' {
								pos += 2
								break
							}
							pos += s
						}
					}

					if collecting {
						if hasPendingCommentRange {
							if !yield(f.NewCommentRange(pendingKind, pendingPos, pendingEnd, pendingHasTrailingNewLine)) {
								return
							}
						}

						pendingPos = startPos
						pendingEnd = pos
						pendingKind = kind
						pendingHasTrailingNewLine = hasTrailingNewLine
						hasPendingCommentRange = true
					}

					continue
				}
				break scan
			default:
				if ch > unicode.MaxASCII && (stringutil.IsWhiteSpaceLike(ch)) {
					if hasPendingCommentRange && stringutil.IsLineBreak(ch) {
						pendingHasTrailingNewLine = true
					}
					pos += size
					continue
				}
				break scan
			}
		}

		if hasPendingCommentRange {
			yield(f.NewCommentRange(pendingKind, pendingPos, pendingEnd, pendingHasTrailingNewLine))
		}
	}
}
