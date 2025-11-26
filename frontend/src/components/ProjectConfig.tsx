import type { ClangFormatStyle } from '../types';

interface ProjectConfigProps {
  projectName: string;
  onProjectNameChange: (name: string) => void;
  cppStandard: number;
  onCppStandardChange: (standard: number) => void;
  includeTests: boolean;
  onIncludeTestsChange: (include: boolean) => void;
  clangFormatStyle: ClangFormatStyle;
  onClangFormatStyleChange: (style: ClangFormatStyle) => void;
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

export function ProjectConfig({
  projectName,
  onProjectNameChange,
  cppStandard,
  onCppStandardChange,
  includeTests,
  onIncludeTestsChange,
  clangFormatStyle,
  onClangFormatStyleChange,
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

        <div className="flex items-center justify-between pt-2">
          <label className="text-sm font-medium text-gray-400">
            Include Test Configuration
          </label>
          <button
            onClick={() => onIncludeTestsChange(!includeTests)}
            className={`relative w-12 h-6 rounded-full transition-all ${
              includeTests ? 'bg-cyan-500' : 'bg-white/10'
            }`}
          >
            <span
              className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-all ${
                includeTests ? 'left-7' : 'left-1'
              }`}
            />
          </button>
        </div>
      </div>
    </div>
  );
}
