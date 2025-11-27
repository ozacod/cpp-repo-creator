import { useState, useEffect } from 'react';
import { fetchVersion, type VersionInfo } from '../api';

interface Architecture {
  id: string;
  name: string;
  filename: string;
  dmgFilename?: string;
}

const MACOS_ARCHS: Architecture[] = [
  { id: 'darwin-arm64', name: 'Apple Silicon', filename: 'forge-darwin-arm64', dmgFilename: 'forge-darwin-arm64.dmg' },
  { id: 'darwin-amd64', name: 'Intel', filename: 'forge-darwin-amd64', dmgFilename: 'forge-darwin-amd64.dmg' },
];

const WINDOWS_ARCHS: Architecture[] = [
  { id: 'windows-amd64', name: 'x86_64', filename: 'forge-windows-amd64.exe' },
];

const LINUX_ARCHS: Architecture[] = [
  { id: 'linux-amd64', name: 'x86_64', filename: 'forge-linux-amd64' },
  { id: 'linux-arm64', name: 'ARM64', filename: 'forge-linux-arm64' },
];

export function CLIDownload() {
  const [versionInfo, setVersionInfo] = useState<VersionInfo>({
    version: '1.0.2',
    cli_version: '1.0.2',
    name: 'forge',
    description: 'C++ Project Generator',
  });
  
  const [macArch, setMacArch] = useState(MACOS_ARCHS[0]);
  const [winArch, setWinArch] = useState(WINDOWS_ARCHS[0]);
  const [linuxArch, setLinuxArch] = useState(LINUX_ARCHS[0]);
  const [copied, setCopied] = useState(false);
  const [showMacDropdown, setShowMacDropdown] = useState(false);
  const [showLinuxDropdown, setShowLinuxDropdown] = useState(false);

  useEffect(() => {
    fetchVersion().then(setVersionInfo).catch(() => {});
  }, []);

  const INSTALL_SCRIPT = 'curl -f https://raw.githubusercontent.com/ozacod/forge/master/install.sh | sh';
  const BASE_URL = 'https://github.com/ozacod/forge/releases/latest/download';

  const copyCommand = () => {
    navigator.clipboard.writeText(INSTALL_SCRIPT);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const formatDate = () => {
    return new Date().toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex flex-col lg:flex-row gap-12 lg:gap-16 items-start">
        
        {/* Left: Logo */}
        <div className="flex-shrink-0">
          <div className="w-48 h-48 lg:w-56 lg:h-56 rounded-3xl bg-gray-900/80 border border-white/10 flex items-center justify-center shadow-2xl">
            <img src="/forge.svg" alt="Forge" className="w-32 h-32 lg:w-40 lg:h-40" />
          </div>
        </div>
        
        {/* Right: Download Options */}
        <div className="flex-1 space-y-8">
          
          {/* Version Info */}
          <div>
            <h1 className="text-5xl lg:text-6xl font-light text-cyan-400 font-mono mb-2">
              {versionInfo.cli_version}
        </h1>
            <div className="flex items-center gap-6 text-gray-400 font-mono text-sm">
              <span>{formatDate()}</span>
              <a 
                href="https://github.com/ozacod/forge/releases" 
            target="_blank" 
            rel="noopener noreferrer"
                className="text-cyan-400 hover:text-cyan-300 flex items-center gap-1"
          >
                View changelog →
          </a>
        </div>
      </div>

          {/* macOS */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                </svg>
                <span className="text-white">macOS</span>
        </div>
              <div className="relative">
                <button 
                  onClick={() => setShowMacDropdown(!showMacDropdown)}
                  className="flex items-center gap-2 text-gray-400 hover:text-white text-sm"
                >
                  {macArch.name}
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                {showMacDropdown && (
                  <div className="absolute right-0 mt-2 w-40 bg-gray-900 border border-white/10 rounded-lg shadow-xl z-10">
                    {MACOS_ARCHS.map((arch) => (
          <button
                        key={arch.id}
                        onClick={() => { setMacArch(arch); setShowMacDropdown(false); }}
                        className={`block w-full text-left px-4 py-2 text-sm hover:bg-white/5 ${
                          macArch.id === arch.id ? 'text-cyan-400' : 'text-gray-300'
                        }`}
                      >
                        {arch.name}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <a
              href={`${BASE_URL}/${macArch.dmgFilename}`}
              className="flex items-center justify-center gap-2 w-full py-3 bg-cyan-500 hover:bg-cyan-400 text-black font-semibold rounded-lg transition-colors"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              Download DMG
            </a>
            <p className="text-xs text-gray-500">
              Requires macOS 10.15 or later •{' '}
              <a 
                href={`${BASE_URL}/${macArch.filename}`}
                className="text-cyan-400 hover:text-cyan-300"
              >
                Download binary only →
              </a>
            </p>
            <details className="text-xs text-gray-500 mt-2">
              <summary className="cursor-pointer hover:text-gray-400">
                ⚠️ "Cannot verify developer" warning? Click here
              </summary>
              <div className="mt-2 p-3 bg-gray-900/80 rounded-lg border border-white/10 space-y-2">
                <p><strong className="text-white">Option 1:</strong> Right-click the app → Open → Click "Open" in dialog</p>
                <p><strong className="text-white">Option 2:</strong> Run in Terminal:</p>
                <code className="block bg-black/50 px-2 py-1 rounded text-cyan-400 font-mono">
                  xattr -d com.apple.quarantine ~/Downloads/forge
                </code>
              </div>
            </details>
          </div>

          {/* Windows */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
                </svg>
                <span className="text-white">Windows</span>
              </div>
              <span className="text-gray-400 text-sm">x86_64</span>
            </div>
            <a
              href={`${BASE_URL}/${winArch.filename}`}
              className="flex items-center justify-center gap-2 w-full py-3 bg-cyan-500 hover:bg-cyan-400 text-black font-semibold rounded-lg transition-colors"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              Download Now
            </a>
          </div>

          {/* Linux */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139z"/>
                </svg>
                <span className="text-white">Linux</span>
              </div>
              <div className="relative">
                <button 
                  onClick={() => setShowLinuxDropdown(!showLinuxDropdown)}
                  className="flex items-center gap-2 text-gray-400 hover:text-white text-sm"
                >
                  {linuxArch.name}
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
          </button>
                {showLinuxDropdown && (
                  <div className="absolute right-0 mt-2 w-32 bg-gray-900 border border-white/10 rounded-lg shadow-xl z-10">
                    {LINUX_ARCHS.map((arch) => (
          <button
                        key={arch.id}
                        onClick={() => { setLinuxArch(arch); setShowLinuxDropdown(false); }}
                        className={`block w-full text-left px-4 py-2 text-sm hover:bg-white/5 ${
                          linuxArch.id === arch.id ? 'text-cyan-400' : 'text-gray-300'
                        }`}
                      >
                        {arch.name}
          </button>
                    ))}
                  </div>
                )}
              </div>
        </div>

            {/* Linux install command */}
            <div className="flex items-center gap-2 bg-gray-900/80 border border-white/10 rounded-lg px-4 py-3">
              <code className="flex-1 text-gray-300 font-mono text-sm overflow-x-auto">
                {INSTALL_SCRIPT}
            </code>
          <button
                onClick={copyCommand}
                className="text-gray-400 hover:text-white p-1"
            title="Copy to clipboard"
          >
                {copied ? (
              <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            ) : (
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            )}
          </button>
        </div>

            <p className="text-xs text-gray-500">
              You can also{' '}
              <a 
                href={`${BASE_URL}/${linuxArch.filename}`}
                className="text-cyan-400 hover:text-cyan-300"
              >
                download the binary directly →
              </a>
        </p>
      </div>

          {/* Separator */}
          <hr className="border-white/10" />

          {/* GitHub link */}
          <p className="text-sm text-gray-500">
            By downloading and using Forge, you agree to its{' '}
            <a 
              href="https://github.com/ozacod/forge/blob/master/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
              className="text-cyan-400 hover:text-cyan-300"
            >
              MIT license
            </a>
            .
          </p>
        </div>
      </div>
    </div>
  );
}
