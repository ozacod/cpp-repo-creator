import type { Library } from '../types';

interface SelectedLibrariesProps {
  libraries: Library[];
  selectedIds: string[];
  onRemove: (id: string) => void;
}

export function SelectedLibraries({ libraries, selectedIds, onRemove }: SelectedLibrariesProps) {
  const selectedLibraries = libraries.filter((lib) => selectedIds.includes(lib.id));

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
              {libs.map((lib) => (
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
                  </div>
                  <button
                    onClick={() => onRemove(lib.id)}
                    className="opacity-0 group-hover:opacity-100 text-gray-500 hover:text-red-400 transition-all"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

