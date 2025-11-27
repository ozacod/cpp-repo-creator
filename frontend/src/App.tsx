import { useState, useEffect, useMemo } from 'react';
import type { Library, Category, LibrarySelection, ProjectConfig, ClangFormatStyle, TestingFramework, ProjectType } from './types';
import { fetchLibraries, fetchCategories, generateProject, previewCMake } from './api';
import { LibraryCard } from './components/LibraryCard';
import { CategoryFilter } from './components/CategoryFilter';
import { SearchBar } from './components/SearchBar';
import { ProjectConfig as ProjectConfigPanel } from './components/ProjectConfig';
import { OptionsModal } from './components/OptionsModal';
import { CLIDownload } from './components/CLIDownload';

type Tab = 'libraries' | 'cli';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('cli');
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

  const [generating, setGenerating] = useState(false);
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

  const handleGenerate = async () => {
    if (!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)) {
      return;
    }

    setGenerating(true);
    try {
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

      const blob = await generateProject(config);

      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${projectName}.zip`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate project');
    } finally {
      setGenerating(false);
    }
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
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-3">
                {/* C++ Logo */}
                <svg className="w-10 h-10" viewBox="0 0 306 344.35" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path fill="#00599C" d="M302.107 258.262c2.401-4.159 3.893-8.845 3.893-13.053V99.14c0-4.208-1.49-8.893-3.892-13.052L153 172.175l149.107 86.087z"/>
                  <path fill="#004482" d="m166.25 341.193 126.5-73.034c3.644-2.104 6.956-5.737 9.357-9.897L153 172.175 3.893 258.263c2.401 4.159 5.714 7.793 9.357 9.896l126.5 73.034c7.287 4.208 19.213 4.208 26.5 0z"/>
                  <path fill="#659AD2" d="M302.108 86.087c-2.402-4.16-5.715-7.793-9.358-9.897L166.25 3.156c-7.287-4.208-19.213-4.208-26.5 0L13.25 76.19C5.962 80.397 0 90.725 0 99.14v146.069c0 4.208 1.491 8.894 3.893 13.053L153 172.175l149.108-86.088z"/>
                  <path fill="#fff" d="M153 274.175c-56.243 0-102-45.757-102-102s45.757-102 102-102c36.292 0 70.139 19.53 88.331 50.968l-44.143 25.544c-9.105-15.736-26.038-25.512-44.188-25.512-28.122 0-51 22.878-51 51 0 28.121 22.878 51 51 51 18.152 0 35.085-9.776 44.191-25.515l44.143 25.543c-18.192 31.441-52.04 50.972-88.334 50.972z"/>
                  <path fill="#fff" d="M255 166.508h-11.334v-11.333h-11.332v11.333H221v11.333h11.334v11.334h11.332v-11.334H255zM297.5 166.508h-11.334v-11.333h-11.332v11.333H263.5v11.333h11.334v11.334h11.332v-11.334H297.5z"/>
                </svg>
                <div>
                  <h1 className="font-display font-bold text-xl text-white">forge</h1>
                  <p className="text-xs text-gray-500">Build modern C++ projects with ease</p>
                </div>
              </div>

              {/* Tabs */}
              <nav className="flex items-center gap-1 bg-black/30 rounded-xl p-1">
                <button
                  onClick={() => setActiveTab('cli')}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                    activeTab === 'cli'
                      ? 'bg-white/10 text-white'
                      : 'text-gray-400 hover:text-white hover:bg-white/5'
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    CLI Tool
                  </span>
                </button>
                <button
                  onClick={() => setActiveTab('libraries')}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                    activeTab === 'libraries'
                      ? 'bg-white/10 text-white'
                      : 'text-gray-400 hover:text-white hover:bg-white/5'
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                    </svg>
                    Libraries
                  </span>
                </button>
              </nav>
            </div>

            {/* Actions - only show for libraries tab */}
            {activeTab === 'libraries' && (
              <div className="flex items-center gap-4">
                <label className="flex items-center gap-2 text-sm text-gray-400">
                  <input
                    type="checkbox"
                    checked={buildShared}
                    onChange={(e) => setBuildShared(e.target.checked)}
                    className="checkbox-custom w-4 h-4"
                  />
                  Shared Libs
                </label>
                
                <button
                  onClick={handleGenerate}
                  disabled={generating || !projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)}
                  className="btn-primary px-6 py-3 rounded-xl font-semibold text-white flex items-center gap-2"
                >
                  {generating ? (
                    <>
                      <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      Generating...
                    </>
                  ) : (
                    <>
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                      Download ZIP
                    </>
                  )}
                </button>
              </div>
            )}
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
        {activeTab === 'cli' ? (
          <CLIDownload />
        ) : (
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

      {/* Footer */}
      <footer className="border-t border-white/5 bg-black/20 mt-16">
        <div className="max-w-7xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between text-sm text-gray-500">
            <p>Built with FetchContent • Modern CMake 3.20+ • Recipe-based configuration</p>
            <div className="flex items-center gap-4">
              <a href="https://cmake.org/cmake/help/latest/module/FetchContent.html" target="_blank" rel="noopener noreferrer" className="hover:text-cyan-400 transition-colors">
                CMake Docs
              </a>
              <a href="https://github.com" target="_blank" rel="noopener noreferrer" className="hover:text-cyan-400 transition-colors">
                GitHub
              </a>
            </div>
          </div>
        </div>
      </footer>

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
