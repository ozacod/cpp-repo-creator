import { useState, useEffect, useMemo } from 'react';
import type { Library, Category, LibrarySelection, ProjectConfig, ClangFormatStyle, TestingFramework, ProjectType } from './types';
import { fetchLibraries, fetchCategories, previewCMake } from './api';
import { LibraryCard } from './components/LibraryCard';
import { CategoryFilter } from './components/CategoryFilter';
import { SearchBar } from './components/SearchBar';
import { ProjectConfig as ProjectConfigPanel } from './components/ProjectConfig';
import { OptionsModal } from './components/OptionsModal';
import { CLIDownload } from './components/CLIDownload';

type Tab = 'home' | 'cli' | 'libraries' | 'docs';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('home');
  const [libraries, setLibraries] = useState<Library[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [selections, setSelections] = useState<Map<string, LibrarySelection>>(new Map());

  const [projectName, setProjectName] = useState('my_project');
  const [cppStandard, setCppStandard] = useState(17);
  const [includeTests, setIncludeTests] = useState(true);
  const [testingFramework, setTestingFramework] = useState<TestingFramework>('googletest');
  const [buildShared, setBuildShared] = useState(false);
  const [clangFormatStyle, setClangFormatStyle] = useState<ClangFormatStyle>('Google');
  const [projectType, setProjectType] = useState<ProjectType>('exe');

  const [optionsLibrary, setOptionsLibrary] = useState<Library | null>(null);

  const [cmakePreview, setCmakePreview] = useState<string>('');

  useEffect(() => {
    Promise.all([fetchLibraries(), fetchCategories()])
      .then(([libs, cats]) => {
        setLibraries(libs);
        setCategories(cats);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  // Update CMake preview when config changes
  useEffect(() => {
    if (!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)) {
      setCmakePreview('# Enter a valid project name to see the preview');
      return;
    }

    const config: ProjectConfig = {
      project_name: projectName,
      cpp_standard: cppStandard,
      libraries: Array.from(selections.values()),
      include_tests: includeTests,
      testing_framework: testingFramework,
      build_shared: buildShared,
      clang_format_style: clangFormatStyle,
      project_type: projectType,
    };

    previewCMake(config)
      .then(setCmakePreview)
      .catch(() => setCmakePreview('# Error generating preview'));
  }, [projectName, cppStandard, selections, includeTests, testingFramework, buildShared, clangFormatStyle, projectType]);

  const filteredLibraries = useMemo(() => {
    let result = libraries;

    if (selectedCategory) {
      result = result.filter((lib) => lib.category === selectedCategory);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter(
        (lib) =>
          lib.name.toLowerCase().includes(query) ||
          lib.description.toLowerCase().includes(query) ||
          lib.tags.some((tag) => tag.toLowerCase().includes(query))
      );
    }

    return result;
  }, [libraries, selectedCategory, searchQuery]);

  const toggleLibrary = (id: string) => {
    setSelections((prev) => {
      const newSelections = new Map(prev);
      if (newSelections.has(id)) {
        newSelections.delete(id);
      } else {
        newSelections.set(id, { library_id: id, options: {} });
      }
      return newSelections;
    });
  };

  const updateLibraryOptions = (libraryId: string, options: Record<string, any>) => {
    setSelections((prev) => {
      const newSelections = new Map(prev);
      const existing = newSelections.get(libraryId);
      if (existing) {
        newSelections.set(libraryId, { ...existing, options });
      } else {
        newSelections.set(libraryId, { library_id: libraryId, options });
      }
      return newSelections;
    });
  };

  const handleConfigureOptions = (library: Library) => {
    // Auto-select library if not already selected
    if (!selections.has(library.id)) {
      setSelections((prev) => {
        const newSelections = new Map(prev);
        newSelections.set(library.id, { library_id: library.id, options: {} });
        return newSelections;
      });
    }
    setOptionsLibrary(library);
  };

  const handleSaveOptions = (options: Record<string, any>) => {
    if (optionsLibrary) {
      updateLibraryOptions(optionsLibrary.id, options);
    }
    setOptionsLibrary(null);
  };

  const generateForgeYaml = () => {
    const selectedLibs = Array.from(selections.values());
    
    // Build dependencies object
    const dependencies: Record<string, Record<string, any>> = {};
    selectedLibs.forEach(sel => {
      const lib = libraries.find(l => l.id === sel.library_id);
      if (lib) {
        dependencies[lib.id] = Object.keys(sel.options).length > 0 ? sel.options : {};
      }
    });

    // Build YAML content
    let yaml = `# forge.yaml - Generated by Forge Web UI
# Run 'forge generate' to create the project

package:
  name: "${projectName}"
  version: "0.1.0"
  cpp_standard: ${cppStandard}

build:
  shared_libs: ${buildShared}
  clang_format: "${clangFormatStyle}"
`;

    if (includeTests && testingFramework !== 'none') {
      yaml += `
testing:
  framework: "${testingFramework}"
`;
    }

    if (Object.keys(dependencies).length > 0) {
      yaml += `
dependencies:
`;
      Object.entries(dependencies).forEach(([libId, options]) => {
        if (Object.keys(options).length > 0) {
          yaml += `  ${libId}:\n`;
          Object.entries(options).forEach(([key, value]) => {
            yaml += `    ${key}: ${value}\n`;
          });
        } else {
          yaml += `  ${libId}: {}\n`;
        }
      });
    }

    return yaml;
  };

  const handleExportYaml = () => {
    if (!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)) {
      return;
    }

    const yaml = generateForgeYaml();
    const blob = new Blob([yaml], { type: 'text/yaml' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
    a.download = 'forge.yaml';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
  };

  if (loading) {
    return (
      <div className="min-h-screen grid-pattern flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-cyan-400/30 border-t-cyan-400 rounded-full animate-spin mx-auto mb-4" />
          <p className="text-gray-400 font-mono">Loading recipes...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen grid-pattern">
      {/* Header */}
      <header className="border-b border-white/5 bg-black/20 backdrop-blur-sm sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-6 py-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              {/* Logo - clickable to go home */}
              <button 
                onClick={() => setActiveTab('home')}
                className="hover:opacity-80 transition-opacity"
              >
                <img src="/forge.svg" alt="Forge" className="w-8 h-8" />
              </button>

              {/* Tabs */}
              <nav className="flex items-center gap-1">
                <button
                  onClick={() => setActiveTab('libraries')}
                  className={`px-3 py-1 rounded text-sm font-medium transition-all ${
                    activeTab === 'libraries'
                      ? 'text-white'
                      : 'text-gray-400 hover:text-white'
                  }`}
                >
                    Libraries
                </button>
                <button
                  onClick={() => setActiveTab('docs')}
                  className={`px-3 py-1 rounded text-sm font-medium transition-all ${
                    activeTab === 'docs'
                      ? 'text-white'
                      : 'text-gray-400 hover:text-white'
                  }`}
                >
                  Docs
                </button>
              </nav>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-3">
              {/* Libraries tab specific actions */}
            {activeTab === 'libraries' && (
                <>
                  <label className="flex items-center gap-2 text-xs text-gray-400">
                  <input
                    type="checkbox"
                    checked={buildShared}
                    onChange={(e) => setBuildShared(e.target.checked)}
                      className="checkbox-custom w-3.5 h-3.5"
                  />
                  Shared Libs
                </label>
                
                <button
                    onClick={handleExportYaml}
                    disabled={!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)}
                    className="btn-primary px-4 py-1.5 rounded-lg text-sm font-semibold text-white flex items-center gap-2"
                  >
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    Export forge.yaml
                  </button>
                    </>
                  )}

              {/* GitHub link */}
              <a
                href="https://github.com/ozacod/forge"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-400 hover:text-white p-2 transition-colors"
                title="View on GitHub"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
                </svg>
              </a>

              {/* Download button - always visible, goes to CLI tab */}
              <button
                onClick={() => setActiveTab('cli')}
                className="bg-cyan-500 hover:bg-cyan-400 text-black px-4 py-1.5 rounded-lg text-sm font-semibold flex items-center gap-2 transition-colors"
              >
                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Download
                <span className="text-xs bg-black/20 px-1 py-0.5 rounded font-mono">D</span>
                </button>
              </div>
          </div>
        </div>
      </header>

      {/* Error banner */}
      {error && (
        <div className="bg-red-500/10 border-b border-red-500/20 px-6 py-3">
          <div className="max-w-7xl mx-auto flex items-center justify-between">
            <span className="text-red-400 text-sm">{error}</span>
            <button onClick={() => setError(null)} className="text-red-400 hover:text-white">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {activeTab === 'home' && (
          <div className="flex flex-col items-center justify-center min-h-[70vh] animate-fade-in">
            <img src="/forge.svg" alt="Forge" className="w-32 h-32 mb-8" />
            <h1 className="font-display text-5xl font-bold text-white mb-4">forge</h1>
            <p className="text-xl text-gray-400 mb-12 text-center max-w-xl">
              The modern C++ project generator. Build, configure, and manage your C++ projects with ease.
            </p>
            <div className="flex gap-4">
              <button
                onClick={() => setActiveTab('cli')}
                className="bg-cyan-500 hover:bg-cyan-400 text-black px-8 py-3 rounded-lg font-semibold flex items-center gap-2 transition-colors"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Download CLI
              </button>
              <button
                onClick={() => setActiveTab('libraries')}
                className="bg-white/10 hover:bg-white/20 text-white px-8 py-3 rounded-lg font-semibold flex items-center gap-2 transition-colors border border-white/10"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                Browse Libraries
              </button>
            </div>
          </div>
        )}

        {activeTab === 'cli' && <CLIDownload />}

        {activeTab === 'docs' && (
          <div className="animate-fade-in max-w-5xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Quick Start */}
              <div className="card-glass rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <span className="text-cyan-400">⚡</span> Quick Start
                </h3>
                <div className="space-y-3 font-mono text-sm">
                  <div className="bg-black/40 rounded-lg p-3">
                    <p className="text-gray-500 text-xs mb-1"># Install forge</p>
                    <code className="text-cyan-400">curl -f https://raw.githubusercontent.com/ozacod/forge/master/install.sh | sh</code>
                  </div>
                  <div className="bg-black/40 rounded-lg p-3">
                    <p className="text-gray-500 text-xs mb-1"># Create new project</p>
                    <code className="text-cyan-400">forge new my_project</code>
                  </div>
                  <div className="bg-black/40 rounded-lg p-3">
                    <p className="text-gray-500 text-xs mb-1"># Generate & build</p>
                    <code className="text-cyan-400">forge generate && forge build</code>
                  </div>
                </div>
              </div>

              {/* Project Commands */}
              <div className="card-glass rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <span className="text-cyan-400">📦</span> Project Commands
                </h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge new [name]</code>
                    <span className="text-gray-400">Create project</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge new --lib</code>
                    <span className="text-gray-400">Create library</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge generate</code>
                    <span className="text-gray-400">Generate CMake files</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge build</code>
                    <span className="text-gray-400">Compile project</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <code className="text-cyan-400">forge run</code>
                    <span className="text-gray-400">Build & run</span>
                  </div>
                </div>
              </div>

              {/* Dependency Commands */}
              <div className="card-glass rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <span className="text-cyan-400">📚</span> Dependencies
                </h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge add spdlog</code>
                    <span className="text-gray-400">Add library</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge remove fmt</code>
                    <span className="text-gray-400">Remove library</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge list</code>
                    <span className="text-gray-400">List available</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <code className="text-cyan-400">forge search json</code>
                    <span className="text-gray-400">Search libraries</span>
                  </div>
                </div>
              </div>

              {/* Other Commands */}
              <div className="card-glass rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <span className="text-cyan-400">🛠️</span> Other Commands
                </h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge test</code>
                    <span className="text-gray-400">Run tests</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge clean</code>
                    <span className="text-gray-400">Clean build</span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-white/5">
                    <code className="text-cyan-400">forge fmt</code>
                    <span className="text-gray-400">Format code</span>
                  </div>
                  <div className="flex justify-between py-1">
                    <code className="text-cyan-400">forge upgrade</code>
                    <span className="text-gray-400">Update forge</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Project Structure */}
            <div className="card-glass rounded-xl p-6 mt-6">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <span className="text-cyan-400">📁</span> Project Structure
              </h3>
              <div className="bg-black/40 rounded-lg p-4 font-mono text-sm text-gray-300">
                <pre>{`my_project/
├── .cmake/forge/
│   └── dependencies.cmake   # Auto-managed by forge add/remove
├── CMakeLists.txt           # Main build file
├── forge.yaml               # Project configuration
├── include/my_project/
│   └── my_project.hpp
├── src/
│   ├── main.cpp
│   └── my_project.cpp
└── tests/`}</pre>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'libraries' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Left column - Library browser */}
            <div className="lg:col-span-2 space-y-6">
              {/* Search and filters */}
              <div className="space-y-4 animate-fade-in">
                <SearchBar value={searchQuery} onChange={setSearchQuery} />
                <CategoryFilter
                  categories={categories}
                  selectedCategory={selectedCategory}
                  onSelect={setSelectedCategory}
                />
              </div>

              {/* Library grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {filteredLibraries.map((lib, index) => (
                  <div
                    key={lib.id}
                    className="animate-slide-up"
                    style={{ animationDelay: `${Math.min(index * 50, 500)}ms` }}
                  >
                    <LibraryCard
                      library={lib}
                      selection={selections.get(lib.id)}
                      onToggle={toggleLibrary}
                      onConfigureOptions={handleConfigureOptions}
                    />
                  </div>
                ))}
              </div>

              {filteredLibraries.length === 0 && (
                <div className="text-center py-16">
                  <svg className="w-16 h-16 mx-auto mb-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p className="text-gray-500">No libraries found</p>
                  <p className="text-sm text-gray-600 mt-1">Try a different search or category</p>
                </div>
              )}

              {/* CMake Preview */}
              <div className="animate-fade-in" style={{ animationDelay: '200ms' }}>
                <CMakePreviewPanel content={cmakePreview} />
              </div>
            </div>

            {/* Right column - Configuration */}
            <div className="space-y-6">
              <div className="animate-slide-up stagger-1">
                <ProjectConfigPanel
                  projectName={projectName}
                  onProjectNameChange={setProjectName}
                  cppStandard={cppStandard}
                  onCppStandardChange={setCppStandard}
                  includeTests={includeTests}
                  onIncludeTestsChange={setIncludeTests}
                  testingFramework={testingFramework}
                  onTestingFrameworkChange={setTestingFramework}
                  clangFormatStyle={clangFormatStyle}
                  onClangFormatStyleChange={setClangFormatStyle}
                  projectType={projectType}
                  onProjectTypeChange={setProjectType}
                />
              </div>

              <div className="animate-slide-up stagger-2">
                <SelectedLibrariesPanel
                  libraries={libraries}
                  selections={selections}
                  onRemove={toggleLibrary}
                  onConfigureOptions={handleConfigureOptions}
                />
              </div>

              {/* Quick tips */}
              <div className="card-glass rounded-2xl p-5 animate-slide-up stagger-3">
                <h3 className="font-display font-semibold text-white mb-3 flex items-center gap-2">
                  <svg className="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Quick Tips
                </h3>
                <ul className="text-sm text-gray-400 space-y-2">
                  <li className="flex items-start gap-2">
                    <span className="text-cyan-400 mt-1">•</span>
                    <span>Click the ⚙️ Options button to configure library-specific build options</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-cyan-400 mt-1">•</span>
                    <span>Recipes are loaded from YAML files - you can add your own!</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-cyan-400 mt-1">•</span>
                    <span>The C++ standard will auto-adjust to the highest required</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-cyan-400 mt-1">•</span>
                    <span>Try the CLI tab for a terminal-based workflow!</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        )}
      </main>


      {/* Options Modal */}
      {optionsLibrary && (
        <OptionsModal
          library={optionsLibrary}
          currentOptions={selections.get(optionsLibrary.id)?.options || {}}
          onSave={handleSaveOptions}
          onClose={() => setOptionsLibrary(null)}
        />
      )}
    </div>
  );
}

// CMake Preview Panel Component
function CMakePreviewPanel({ content }: { content: string }) {
  const [expanded, setExpanded] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(content);
  };

  return (
    <div className="card-glass rounded-2xl overflow-hidden">
      <div className="flex items-center justify-between px-5 py-3 border-b border-white/10">
        <div className="flex items-center gap-2">
          <span className="font-mono text-sm text-cyan-400">CMakeLists.txt</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={copyToClipboard}
            className="p-2 text-gray-500 hover:text-white transition-colors"
            title="Copy to clipboard"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </button>
          <button
            onClick={() => setExpanded(!expanded)}
            className="p-2 text-gray-500 hover:text-white transition-colors"
            title={expanded ? 'Collapse' : 'Expand'}
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {expanded ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              )}
            </svg>
          </button>
        </div>
      </div>
      <pre
        className={`code-preview p-5 overflow-x-auto text-sm transition-all ${
          expanded ? 'max-h-[800px]' : 'max-h-[300px]'
        }`}
      >
        <code className="text-gray-300">
          {content.split('\n').map((line, i) => (
            <div key={i} className="flex">
              <span className="w-8 text-right pr-4 text-gray-600 select-none">{i + 1}</span>
              <span className={
                line.startsWith('#') ? 'text-green-400' :
                line.includes('FetchContent') ? 'text-purple-400' :
                line.includes('target_') ? 'text-cyan-400' :
                line.includes('set(') || line.includes('add_') ? 'text-yellow-400' :
                ''
              }>
                {line || ' '}
              </span>
            </div>
          ))}
        </code>
      </pre>
    </div>
  );
}

// Selected Libraries Panel Component
function SelectedLibrariesPanel({
  libraries,
  selections,
  onRemove,
  onConfigureOptions,
}: {
  libraries: Library[];
  selections: Map<string, LibrarySelection>;
  onRemove: (id: string) => void;
  onConfigureOptions: (library: Library) => void;
}) {
  const selectedLibraries = libraries.filter((lib) => selections.has(lib.id));

  if (selectedLibraries.length === 0) {
    return (
      <div className="card-glass rounded-2xl p-6">
        <h2 className="font-display font-semibold text-lg text-white flex items-center gap-2 mb-4">
          <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          Selected Libraries
        </h2>
        <div className="text-center py-8 text-gray-500">
          <svg className="w-12 h-12 mx-auto mb-3 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <p className="text-sm">No libraries selected</p>
          <p className="text-xs mt-1">Click on libraries to add them</p>
        </div>
      </div>
    );
  }

  // Group by category
  const grouped = selectedLibraries.reduce((acc, lib) => {
    if (!acc[lib.category]) acc[lib.category] = [];
    acc[lib.category].push(lib);
    return acc;
  }, {} as Record<string, Library[]>);

  return (
    <div className="card-glass rounded-2xl p-6">
      <h2 className="font-display font-semibold text-lg text-white flex items-center gap-2 mb-4">
        <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Selected Libraries
        <span className="ml-auto text-sm font-mono text-cyan-400">
          {selectedLibraries.length}
        </span>
      </h2>

      <div className="space-y-4 max-h-[400px] overflow-y-auto pr-2">
        {Object.entries(grouped).map(([category, libs]) => (
          <div key={category}>
            <h3 className="text-xs font-mono uppercase tracking-wider text-gray-500 mb-2">
              {category}
            </h3>
            <div className="space-y-2">
              {libs.map((lib) => {
                const selection = selections.get(lib.id);
                const optionsCount = selection ? Object.keys(selection.options).length : 0;
                
                return (
                  <div
                    key={lib.id}
                    className="flex items-center justify-between bg-white/5 rounded-lg px-3 py-2 group"
                  >
                    <div className="flex items-center gap-2 min-w-0">
                      <span className="w-2 h-2 rounded-full bg-cyan-400" />
                      <span className="text-sm text-white truncate">{lib.name}</span>
                      <span className="text-[10px] font-mono text-gray-500">
                        C++{lib.cpp_standard}
                      </span>
                      {optionsCount > 0 && (
                        <span className="px-1.5 py-0.5 text-[10px] bg-purple-500/20 text-purple-400 rounded-full">
                          {optionsCount} opt{optionsCount > 1 ? 's' : ''}
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-1">
                      {lib.options.length > 0 && (
                        <button
                          onClick={() => onConfigureOptions(lib)}
                          className="opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-purple-400 transition-all"
                          title="Configure options"
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                          </svg>
                        </button>
                      )}
                      <button
                        onClick={() => onRemove(lib.id)}
                        className="opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-red-400 transition-all"
                        title="Remove"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
        </button>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
