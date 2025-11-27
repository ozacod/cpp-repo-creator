import type { ClangFormatStyle, TestingFramework, ProjectType } from '../types';

interface ProjectConfigProps {
  projectName: string;
  onProjectNameChange: (name: string) => void;
  cppStandard: number;
  onCppStandardChange: (standard: number) => void;
  includeTests: boolean;
  onIncludeTestsChange: (include: boolean) => void;
  testingFramework: TestingFramework;
  onTestingFrameworkChange: (framework: TestingFramework) => void;
  clangFormatStyle: ClangFormatStyle;
  onClangFormatStyleChange: (style: ClangFormatStyle) => void;
  projectType: ProjectType;
  onProjectTypeChange: (type: ProjectType) => void;
}

const CPP_STANDARDS = [11, 14, 17, 20, 23];

const CLANG_FORMAT_STYLES: { id: ClangFormatStyle; name: string; description: string }[] = [
  { id: 'Google', name: 'Google', description: 'Google C++ style guide' },
  { id: 'LLVM', name: 'LLVM', description: 'LLVM coding standards' },
  { id: 'Chromium', name: 'Chromium', description: 'Chromium project style' },
  { id: 'Mozilla', name: 'Mozilla', description: 'Mozilla coding style' },
  { id: 'WebKit', name: 'WebKit', description: 'WebKit coding style' },
  { id: 'Microsoft', name: 'Microsoft', description: 'Microsoft C++ style' },
  { id: 'GNU', name: 'GNU', description: 'GNU coding standards' },
];

const TESTING_FRAMEWORKS: { id: TestingFramework; name: string; description: string }[] = [
  { id: 'googletest', name: 'GoogleTest', description: 'Google\'s C++ testing framework with mocking' },
  { id: 'catch2', name: 'Catch2', description: 'Modern C++ test framework with BDD support' },
  { id: 'doctest', name: 'doctest', description: 'Fast single-header testing framework' },
  { id: 'none', name: 'None', description: 'No testing framework' },
];

const PROJECT_TYPES: { id: ProjectType; name: string; description: string; icon: string }[] = [
  { id: 'exe', name: 'Executable', description: 'Application with main.cpp', icon: '🚀' },
  { id: 'lib', name: 'Library', description: 'Reusable library', icon: '📦' },
];

export function ProjectConfig({
  projectName,
  onProjectNameChange,
  cppStandard,
  onCppStandardChange,
  includeTests,
  onIncludeTestsChange,
  testingFramework,
  onTestingFrameworkChange,
  clangFormatStyle,
  onClangFormatStyleChange,
  projectType,
  onProjectTypeChange,
}: ProjectConfigProps) {
  const isValidName = /^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName) || projectName === '';

  return (
    <div className="card-glass rounded-2xl p-6 space-y-5">
      <h2 className="font-display font-semibold text-lg text-white flex items-center gap-2">
        <svg className="w-5 h-5 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        Project Configuration
      </h2>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-2">
            Project Type
          </label>
          <div className="grid grid-cols-2 gap-2">
            {PROJECT_TYPES.map((type) => (
              <button
                key={type.id}
                onClick={() => onProjectTypeChange(type.id)}
                title={type.description}
                className={`py-3 px-4 rounded-lg text-sm transition-all flex items-center gap-2 ${
                  projectType === type.id
                    ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40'
                    : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                }`}
              >
                <span>{type.icon}</span>
                <span className="font-medium">{type.name}</span>
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-2">
            Project Name
          </label>
          <input
            type="text"
            value={projectName}
            onChange={(e) => onProjectNameChange(e.target.value)}
            placeholder="my_awesome_project"
            className={`input-field w-full px-4 py-3 rounded-lg text-white font-mono ${
              !isValidName && projectName ? 'border-red-500/50 focus:border-red-500' : ''
            }`}
          />
          {!isValidName && projectName && (
            <p className="mt-1.5 text-xs text-red-400">
              Must start with a letter, only letters, numbers, and underscores
            </p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-2">
            C++ Standard
          </label>
          <div className="flex gap-2">
            {CPP_STANDARDS.map((std) => (
              <button
                key={std}
                onClick={() => onCppStandardChange(std)}
                className={`flex-1 py-2.5 rounded-lg font-mono text-sm transition-all ${
                  cppStandard === std
                    ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40'
                    : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                }`}
              >
                C++{std}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-2">
            Testing Framework
          </label>
          <div className="grid grid-cols-4 gap-2">
            {TESTING_FRAMEWORKS.map((fw) => (
              <button
                key={fw.id}
                onClick={() => {
                  onTestingFrameworkChange(fw.id);
                  // Auto-enable tests when selecting a framework
                  if (fw.id !== 'none' && !includeTests) {
                    onIncludeTestsChange(true);
                  }
                  // Auto-disable tests when selecting none
                  if (fw.id === 'none' && includeTests) {
                    onIncludeTestsChange(false);
                  }
                }}
                title={fw.description}
                className={`py-2 px-2 rounded-lg font-mono text-xs transition-all ${
                  testingFramework === fw.id
                    ? 'bg-green-500/20 text-green-400 border border-green-500/40'
                    : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                }`}
              >
                {fw.name}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-2">
            Clang-Format Style
          </label>
          <div className="grid grid-cols-4 gap-2">
            {CLANG_FORMAT_STYLES.map((style) => (
              <button
                key={style.id}
                onClick={() => onClangFormatStyleChange(style.id)}
                title={style.description}
                className={`py-2 px-2 rounded-lg font-mono text-xs transition-all ${
                  clangFormatStyle === style.id
                    ? 'bg-purple-500/20 text-purple-400 border border-purple-500/40'
                    : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10'
                }`}
              >
                {style.name}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
