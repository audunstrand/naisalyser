# Test Task T5: Local Reader Tests

## Description
Implement comprehensive tests for local repository reading:
- `NewReader()`
- `ReadRepository()`
- `GetFileContent()`
- `FileExists()`
- `detectLanguage()`
- `readFirstLineOfReadme()`

## File Location
`internal/local/reader_test.go`

## Test Cases Required

### NewReader Tests (2 tests)
1. `TestNewReader_VerboseFalse` - Creates reader with verbose=false
2. `TestNewReader_VerboseTrue` - Creates reader with verbose=true

### ReadRepository Tests (15 tests)
1. `TestReadRepository_ValidPath` - Reads valid repository directory
2. `TestReadRepository_AbsolutePathConversion` - Handles relative paths
3. `TestReadRepository_InvalidPath` - Error for non-existent path
4. `TestReadRepository_NotDirectory` - Error for file path instead of directory
5. `TestReadRepository_NameExtraction` - Correctly extracts directory name
6. `TestReadRepository_OrgParameter` - Uses org parameter in FullName
7. `TestReadRepository_FullNameFormat` - Format is "org/name"
8. `TestReadRepository_DefaultOrg` - Works with default "navikt" org
9. `TestReadRepository_CustomOrg` - Works with custom org name
10. `TestReadRepository_LanguageDetection` - Detects language correctly
11. `TestReadRepository_ReadmeExtraction` - Extracts README first line
12. `TestReadRepository_NoReadme` - Handles missing README gracefully
13. `TestReadRepository_EmptyRepository` - Handles directory with no files
14. `TestReadRepository_Symlinks` - Handles symlinked directories
15. `TestReadRepository_PermissionError` - Error for inaccessible directory

### GetFileContent Tests (8 tests)
1. `TestGetFileContent_ExistingFile` - Reads file content correctly
2. `TestGetFileContent_NonExistentFile` - Error for missing file
3. `TestGetFileContent_IsDirectory` - Error for directory path
4. `TestGetFileContent_EmptyFile` - Handles empty files
5. `TestGetFileContent_LargeFile` - Reads large files
6. `TestGetFileContent_BinaryFile` - Handles binary files
7. `TestGetFileContent_PermissionDenied` - Error for unreadable file
8. `TestGetFileContent_RelativePath` - Works with relative paths

### FileExists Tests (6 tests)
1. `TestFileExists_ExistingFile` - Returns true for existing file
2. `TestFileExists_NonExistentFile` - Returns false for missing file
3. `TestFileExists_IsDirectory` - Returns false for directory
4. `TestFileExists_SymlinkedFile` - Returns true for symlink
5. `TestFileExists_RelativePath` - Works with relative paths
6. `TestFileExists_PermissionDenied` - Returns false for inaccessible file

### DetectLanguage Tests (20 tests)

#### Go Detection (3 tests)
1. `TestDetectLanguage_GoWithGoMod` - Detects go.mod
2. `TestDetectLanguage_GoWithGoSum` - Detects go.sum
3. `TestDetectLanguage_GoWithBoth` - Both files present

#### Java/Kotlin Detection (4 tests)
1. `TestDetectLanguage_JavaWithPomXml` - Detects pom.xml
2. `TestDetectLanguage_KotlinWithGradleKts` - Detects build.gradle.kts
3. `TestDetectLanguage_KotlinSourceFiles` - Detects *.kt files
4. `TestDetectLanguage_JavaSourceFiles` - Detects *.java files

#### JavaScript/TypeScript Detection (4 tests)
1. `TestDetectLanguage_JavaScriptWithPackageJson` - Detects package.json
2. `TestDetectLanguage_TypeScriptWithTsConfig` - Detects tsconfig.json
3. `TestDetectLanguage_JavaScriptSourceFiles` - Detects *.js files
4. `TestDetectLanguage_TypeScriptSourceFiles` - Detects *.ts files

#### Python Detection (2 tests)
1. `TestDetectLanguage_PythonWithRequirements` - Detects requirements.txt
2. `TestDetectLanguage_PythonWithSetup` - Detects setup.py

#### Other (4 tests)
1. `TestDetectLanguage_MultipleLanguages` - Returns first detected (priority)
2. `TestDetectLanguage_UnknownLanguage` - Returns empty string
3. `TestDetectLanguage_EmptyDirectory` - Returns empty string
4. `TestDetectLanguage_HiddenFiles` - Ignores hidden files

#### Priority Tests (3 tests)
1. `TestDetectLanguage_JavaPriority` - Java detected before Node
2. `TestDetectLanguage_GoOverride` - Go overrides others
3. `TestDetectLanguage_FirstFound` - Stops at first match

### ReadFirstLineOfReadme Tests (6 tests)
1. `TestReadFirstLineOfReadme_Valid` - Reads first non-empty line
2. `TestReadFirstLineOfReadme_MultipleLines` - Ignores remaining lines
3. `TestReadFirstLineOfReadme_NoReadme` - Returns empty string when missing
4. `TestReadFirstLineOfReadme_EmptyReadme` - Handles empty README
5. `TestReadFirstLineOfReadme_OnlyNewlines` - Handles whitespace-only file
6. `TestReadFirstLineOfReadme_LongFirstLine` - Handles very long description

## Test Repository Structures

### Go Project
```
go-project/
├── go.mod
├── go.sum
├── README.md
└── main.go
```

### Java Project
```
java-project/
├── pom.xml
├── README.md
└── src/
```

### Kotlin Project
```
kotlin-project/
├── build.gradle.kts
├── README.md
└── src/main/kotlin/
```

### Node Project
```
node-project/
├── package.json
├── tsconfig.json
├── README.md
├── src/
└── package-lock.json
```

### Python Project
```
python-project/
├── requirements.txt
├── README.md
└── setup.py
```

## Edge Cases to Test
- Very deep directory structures
- Special characters in directory names
- Unicode in file content
- Very long file paths
- Circular symlinks
- Files with unusual encodings
- Mixed language projects

## Implementation Notes
- Use `t.TempDir()` to create test repos
- Use `t.Helper()` for test helper functions
- Create minimal valid files (no actual code)
- Test both existence checks and content reading
- Handle permission errors gracefully (skip on Windows if needed)

## Success Criteria
- ✅ All 25+ tests pass
- ✅ Code coverage for reader functions ≥ 90%
- ✅ All language detections verified
- ✅ README extraction works with real files
- ✅ Error handling verified
