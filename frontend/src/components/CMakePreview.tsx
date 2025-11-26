import { useState, useEffect } from 'react';
import { previewCMake } from '../api';

interface CMakePreviewProps {
  projectName: string;
  cppStandard: number;
  libraryIds: string[];
  includeTests: boolean;
}

export function CMakePreview({ projectName, cppStandard, libraryIds, includeTests }: CMakePreviewProps) {
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [expanded, setExpanded] = useState(false);

  useEffect(() => {
    if (!projectName || !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(projectName)) {
      setContent('# Enter a valid project name to see the preview');
      return;
    }

    setLoading(true);
    previewCMake(projectName, cppStandard, libraryIds, includeTests)
      .then(setContent)
      .catch(() => setContent('# Error generating preview'))
      .finally(() => setLoading(false));
  }, [projectName, cppStandard, libraryIds, includeTests]);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(content);
  };

  return (
    <div className="card-glass rounded-2xl overflow-hidden">
      <div className="flex items-center justify-between px-5 py-3 border-b border-white/10">
        <div className="flex items-center gap-2">
          <span className="font-mono text-sm text-cyan-400">CMakeLists.txt</span>
          {loading && (
            <span className="w-4 h-4 border-2 border-cyan-400/30 border-t-cyan-400 rounded-full animate-spin" />
          )}
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

