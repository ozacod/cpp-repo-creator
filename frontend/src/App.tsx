import { useState, useEffect, useMemo } from 'react';
import type { Library, Category } from './types';
import { fetchLibraries, fetchCategories, generateProject } from './api';
import { LibraryCard } from './components/LibraryCard';
import { CategoryFilter } from './components/CategoryFilter';
import { SearchBar } from './components/SearchBar';
import { ProjectConfig } from './components/ProjectConfig';
import { SelectedLibraries } from './components/SelectedLibraries';
import { CMakePreview } from './components/CMakePreview';

function App() {
  const [libraries, setLibraries] = useState<Library[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  const [projectName, setProjectName] = useState('my_project');
  const [cppStandard, setCppStandard] = useState(17);
  const [includeTests, setIncludeTests] = useState(true);

  const [generating, setGenerating] = useState(false);

  useEffect(() => {
    Promise.all([fetchLibraries(), fetchCategories()])
      .then(([libs, cats]) => {
        setLibraries(libs);
        setCategories(cats);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

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
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const handleGenerate = async () => {
    if (!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)) {
      return;
    }

    setGenerating(true);
    try {
      const blob = await generateProject({
        project_name: projectName,
        cpp_standard: cppStandard,
        library_ids: selectedIds,
        include_tests: includeTests,
      });

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
          <p className="text-gray-400 font-mono">Loading libraries...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen grid-pattern">
      {/* Header */}
      <header className="border-b border-white/5 bg-black/20 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-cyan-400 to-purple-500 flex items-center justify-center">
                <span className="font-mono font-bold text-white text-lg">C++</span>
              </div>
              <div>
                <h1 className="font-display font-bold text-xl text-white">Project Creator</h1>
                <p className="text-xs text-gray-500">Build modern C++ projects with ease</p>
              </div>
            </div>

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
                    isSelected={selectedIds.includes(lib.id)}
                    onToggle={toggleLibrary}
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
              <CMakePreview
                projectName={projectName}
                cppStandard={cppStandard}
                libraryIds={selectedIds}
                includeTests={includeTests}
              />
            </div>
          </div>

          {/* Right column - Configuration */}
          <div className="space-y-6">
            <div className="animate-slide-up stagger-1">
              <ProjectConfig
                projectName={projectName}
                onProjectNameChange={setProjectName}
                cppStandard={cppStandard}
                onCppStandardChange={setCppStandard}
                includeTests={includeTests}
                onIncludeTestsChange={setIncludeTests}
              />
            </div>

            <div className="animate-slide-up stagger-2">
              <SelectedLibraries
                libraries={libraries}
                selectedIds={selectedIds}
                onRemove={toggleLibrary}
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
                  <span>Libraries with the same color tags are alternatives to each other</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-cyan-400 mt-1">•</span>
                  <span>Header-only libraries don't require separate compilation</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-cyan-400 mt-1">•</span>
                  <span>The C++ standard will auto-adjust to the highest required</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-white/5 bg-black/20 mt-16">
        <div className="max-w-7xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between text-sm text-gray-500">
            <p>Built with FetchContent • Modern CMake 3.20+</p>
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
    </div>
  );
}

export default App;
