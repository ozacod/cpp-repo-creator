import { useState, useEffect } from 'react';
import { fetchVersion, type VersionInfo } from '../api';

interface Platform {
  id: string;
  name: string;
  icon: React.ReactNode;
  architectures: Architecture[];
}

interface Architecture {
  id: string;
  name: string;
  filename: string;
  size?: string;
}

const PLATFORMS: Platform[] = [
  {
    id: 'darwin',
    name: 'macOS',
    icon: (
      <svg className="w-8 h-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
      </svg>
    ),
    architectures: [
      { id: 'darwin-arm64', name: 'Apple Silicon (M1/M2/M3)', filename: 'forge-darwin-arm64' },
      { id: 'darwin-amd64', name: 'Intel (x86_64)', filename: 'forge-darwin-amd64' },
    ],
  },
  {
    id: 'linux',
    name: 'Linux',
    icon: (
      <svg className="w-8 h-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139zm.529 3.405h.013c.213 0 .396.062.584.198.19.135.33.332.438.533.105.259.158.459.166.724 0-.02.006-.04.006-.06v.105a.086.086 0 01-.004-.021l-.004-.024a1.807 1.807 0 01-.15.706.953.953 0 01-.213.335.71.71 0 00-.088-.042c-.104-.045-.198-.064-.284-.133a1.312 1.312 0 00-.22-.066c.05-.06.146-.133.183-.198.053-.128.082-.264.088-.402v-.02a1.21 1.21 0 00-.061-.4c-.045-.134-.101-.2-.183-.333-.084-.066-.167-.132-.267-.132h-.016c-.093 0-.176.03-.262.132a.8.8 0 00-.205.334 1.18 1.18 0 00-.09.4v.019c.002.089.008.179.02.267-.193-.067-.438-.135-.607-.202a1.635 1.635 0 01-.018-.2v-.02a1.772 1.772 0 01.15-.768c.082-.22.232-.406.43-.533a.985.985 0 01.594-.2zm-2.962.059h.036c.142 0 .27.048.399.135.146.129.264.288.344.465.09.199.14.4.153.667v.004c.007.134.006.2-.002.266v.08c-.03.007-.056.018-.083.024-.152.055-.274.135-.393.2.012-.09.013-.18.003-.267v-.015c-.012-.133-.04-.2-.082-.333a.613.613 0 00-.166-.267.248.248 0 00-.183-.064h-.021c-.071.006-.13.04-.186.132a.552.552 0 00-.12.27.944.944 0 00-.023.33v.015c.012.135.037.2.08.334.046.134.098.2.166.268.01.009.02.018.034.024-.07.057-.117.07-.176.136a.304.304 0 01-.131.068 2.62 2.62 0 01-.275-.402 1.772 1.772 0 01-.155-.667 1.759 1.759 0 01.08-.668 1.43 1.43 0 01.283-.535c.128-.133.26-.2.418-.2zm1.37 1.706c.332 0 .733.065 1.216.399.293.2.523.269 1.052.468h.003c.255.136.405.266.478.399v-.131a.571.571 0 01.016.47c-.123.31-.516.643-1.063.842v.002c-.268.135-.501.333-.775.465-.276.135-.588.292-1.012.267a1.139 1.139 0 01-.448-.067 3.566 3.566 0 01-.322-.198c-.195-.135-.363-.332-.612-.465v-.005h-.005c-.4-.246-.616-.512-.686-.71-.07-.268-.005-.47.193-.6.224-.135.38-.271.483-.336.104-.074.143-.102.176-.131h.002v-.003c.169-.202.436-.47.839-.601.139-.036.294-.065.466-.065zm2.8 2.142c.358 1.417 1.196 3.475 1.735 4.473.286.534.855 1.659 1.102 3.024.156-.005.33.018.513.064.646-1.671-.546-3.467-1.089-3.966-.22-.2-.232-.335-.123-.335.59.534 1.365 1.572 1.646 2.757.13.535.16 1.104.021 1.67.067.028.135.06.205.067 1.032.534 1.413.938 1.23 1.537v-.002c-.06-.135-.12-.2-.184-.268-.193-.135-.402-.2-.614-.267-.545.2-1.154.074-1.7-.333-.545-.398-1.025-1.068-1.327-1.936-.316.266-.59.535-.854.668-.106.065-.19.1-.257.2-.128.2-.208.467-.086.998.167.666.042 1.201-.283 1.468-.32.27-.663.2-1.074.068a3.204 3.204 0 01-1.02-.466 8.57 8.57 0 00-1.156-.6 3.59 3.59 0 00-.82-.267 5.58 5.58 0 00-.377-.2 3.84 3.84 0 00-.86-.133c-.034-.134-.074-.267-.104-.4-.04-.2-.08-.399-.106-.466a.556.556 0 01-.053-.2c-.006-.135.073-.2.2-.2.052 0 .106.013.161.04.18.087.377.209.58.267.14.04.28.04.418 0a.604.604 0 00.313-.2c.085-.135.136-.267.17-.4.032-.135.052-.267.05-.399 0-.401-.128-.802-.384-1.136-.168-.2-.371-.4-.636-.6a4.018 4.018 0 00-.646-.333c-.155-.068-.334-.135-.467-.135a1.264 1.264 0 00-.39.068c-.11.043-.2.081-.288.132-.19.133-.36.269-.536.466-.14.135-.287.333-.41.535-.12.2-.225.4-.312.667l-.078.268a2.12 2.12 0 00-.054.263c-.016.133-.017.2-.008.267-.09.069-.2.133-.313.2-.1.067-.217.132-.346.2-.055.023-.11.046-.16.067l-.02.003c-.184.067-.328.135-.437.2a.91.91 0 00-.191.135.356.356 0 00-.104.133c-.04.135.006.2.107.266.106.067.266.133.459.133h.024c.134 0 .28-.034.45-.067.164-.035.336-.098.484-.131.1-.035.175-.067.238-.067.143.006.206.135.239.2.038.134.045.2.035.266.006.066.016.133.038.2.023.067.048.098.1.132.122.067.292.065.502.067h.017c.209 0 .43-.028.642-.067.209-.04.418-.098.623-.135.19-.032.38-.053.57-.053.209 0 .348.02.57.054.197.033.4.067.63.098.226.034.443.034.645-.003.203-.032.388-.098.574-.199.06-.032.118-.065.169-.098.056.067.115.135.168.2.054.065.111.133.164.2.055.067.109.132.165.198a.66.66 0 00.214.133.26.26 0 00.119.035c.056 0 .113-.014.164-.042.052-.028.099-.066.136-.108.036-.04.07-.085.1-.13a2.1 2.1 0 00.166-.362v-.003c.032-.098.054-.199.064-.3.014-.11.019-.2.019-.332v-.003c.004-.133 0-.265-.014-.399a2.758 2.758 0 00-.09-.466c-.024-.066-.034-.133-.048-.2h.002c.04.067.086.133.134.2.048.066.094.134.147.2.053.068.11.134.17.2a.63.63 0 00.205.133c.043.02.088.032.134.032h.011c.033 0 .068-.005.1-.015a.336.336 0 00.094-.04.388.388 0 00.083-.061.418.418 0 00.068-.081.538.538 0 00.08-.21c.017-.1.023-.2.016-.302a1.936 1.936 0 00-.058-.362 2.23 2.23 0 00-.134-.4 3.65 3.65 0 00-.191-.4 6.67 6.67 0 00-.246-.468c-.067-.133-.137-.266-.21-.4l-.03-.055a11.76 11.76 0 00-.166-.308c-.048-.088-.1-.175-.153-.263a3.11 3.11 0 00-.198-.31 1.84 1.84 0 00-.26-.3 1.35 1.35 0 00-.338-.23c-.133-.065-.283-.097-.438-.097h-.027c-.165 0-.326.033-.48.098a1.61 1.61 0 00-.412.265 1.89 1.89 0 00-.332.4c-.099.15-.18.31-.247.477-.067.167-.113.333-.143.5-.033.167-.053.333-.053.5v.117c0 .033-.002.067-.005.1l-.002.028v.008c-.003.134-.003.2-.007.266-.01.135-.016.268-.048.4a2.78 2.78 0 00-.067.4c.01.134.024.198.048.333.024.133.062.267.105.4.044.133.098.266.152.4.055.133.113.199.17.332.058.133.115.267.175.4.064.133.122.267.184.4a19.06 19.06 0 00.187.4c.068.133.136.266.208.4l.112.198c.014.03.03.055.046.083z"/>
      </svg>
    ),
    architectures: [
      { id: 'linux-amd64', name: 'x86_64 (64-bit)', filename: 'forge-linux-amd64' },
      { id: 'linux-arm64', name: 'ARM64 (aarch64)', filename: 'forge-linux-arm64' },
    ],
  },
  {
    id: 'windows',
    name: 'Windows',
    icon: (
      <svg className="w-8 h-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
      </svg>
    ),
    architectures: [
      { id: 'windows-amd64', name: 'x86_64 (64-bit)', filename: 'forge-windows-amd64.exe' },
    ],
  },
];

const FEATURES = [
  {
    icon: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
    ),
    title: 'Fast & Lightweight',
    description: 'Single binary, no dependencies. Just download and run.',
  },
  {
    icon: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    ),
    title: 'YAML Configuration',
    description: 'Define your project in a simple forge.yaml file.',
  },
  {
    icon: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
    ),
    title: 'CLI-First Workflow',
    description: 'Perfect for scripts, CI/CD, and terminal lovers.',
  },
  {
    icon: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
      </svg>
    ),
    title: '60+ Libraries',
    description: 'Access the same curated library collection.',
  },
];

export function CLIDownload() {
  const [selectedPlatform, setSelectedPlatform] = useState<Platform | null>(null);
  const [copiedCommand, setCopiedCommand] = useState<string | null>(null);
  const [installMethod, setInstallMethod] = useState<'script' | 'manual'>('script');
  const [versionInfo, setVersionInfo] = useState<VersionInfo>({
    version: '1.0.2',
    cli_version: '1.0.2',
    name: 'forge',
    description: 'C++ Project Generator',
  });

  useEffect(() => {
    fetchVersion().then(setVersionInfo).catch(() => {
      // Keep default version on error
    });
  }, []);

  const INSTALL_SCRIPT_CURL = 'sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/forge/master/install.sh)"';
  const INSTALL_SCRIPT_WGET = 'sh -c "$(wget -qO- https://raw.githubusercontent.com/ozacod/forge/master/install.sh)"';

  const copyCommand = (command: string) => {
    navigator.clipboard.writeText(command);
    setCopiedCommand(command);
    setTimeout(() => setCopiedCommand(null), 2000);
  };

  const getInstallCommand = (arch: Architecture) => {
    const baseUrl = 'https://github.com/ozacod/forge/releases/latest/download';
    if (arch.id.startsWith('windows')) {
      return `# Download from:\n${baseUrl}/${arch.filename}\n\n# Or using PowerShell:\nInvoke-WebRequest -Uri "${baseUrl}/${arch.filename}" -OutFile "forge.exe"`;
    }
    return `curl -L -o forge ${baseUrl}/${arch.filename} && chmod +x forge && sudo mv forge /usr/local/bin/`;
  };

  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <div className="text-center space-y-6 animate-fade-in">
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-cyan-500/10 to-purple-500/10 border border-cyan-500/20">
          <span className="text-cyan-400 font-mono text-sm">v{versionInfo.cli_version}</span>
          <span className="text-gray-500">•</span>
          <span className="text-gray-400 text-sm">Go • Static Binary</span>
        </div>
        
        {/* Forge Logo */}
        <div className="flex justify-center">
          <svg className="w-24 h-24" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="forgeGradHero" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style={{stopColor:'#22d3ee', stopOpacity:1}} />
                <stop offset="100%" style={{stopColor:'#a855f7', stopOpacity:1}} />
              </linearGradient>
              <linearGradient id="anvilGradHero" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style={{stopColor:'#4b5563', stopOpacity:1}} />
                <stop offset="100%" style={{stopColor:'#1f2937', stopOpacity:1}} />
              </linearGradient>
              <linearGradient id="fireGradHero" x1="0%" y1="100%" x2="0%" y2="0%">
                <stop offset="0%" style={{stopColor:'#f97316', stopOpacity:1}} />
                <stop offset="50%" style={{stopColor:'#eab308', stopOpacity:1}} />
                <stop offset="100%" style={{stopColor:'#fef08a', stopOpacity:1}} />
              </linearGradient>
            </defs>
            <circle cx="32" cy="32" r="30" fill="url(#forgeGradHero)" opacity="0.15"/>
            <circle cx="32" cy="32" r="30" fill="none" stroke="url(#forgeGradHero)" strokeWidth="2"/>
            <path d="M16 42 L48 42 L52 48 L12 48 Z" fill="url(#anvilGradHero)"/>
            <path d="M20 36 L44 36 L48 42 L16 42 Z" fill="#374151"/>
            <path d="M18 32 L46 32 L44 36 L20 36 Z" fill="#4b5563"/>
            <path d="M10 32 L18 32 L18 36 L14 36 Z" fill="#6b7280"/>
            <rect x="30" y="14" width="4" height="20" rx="1" fill="#78716c" transform="rotate(-30 32 24)"/>
            <rect x="24" y="10" width="14" height="8" rx="2" fill="#57534e" transform="rotate(-30 32 14)"/>
            <ellipse cx="32" cy="30" rx="4" ry="6" fill="url(#fireGradHero)" opacity="0.9"/>
            <ellipse cx="28" cy="28" rx="2" ry="3" fill="#fbbf24" opacity="0.7"/>
            <ellipse cx="36" cy="27" rx="2" ry="3" fill="#fbbf24" opacity="0.7"/>
            <text x="32" y="56" fontFamily="Arial, sans-serif" fontSize="10" fontWeight="bold" fill="white" textAnchor="middle">C++</text>
          </svg>
        </div>
        
        <h1 className="font-display text-5xl font-bold text-white">
          forge <span className="text-cyan-400">CLI</span>
        </h1>
        
        <p className="text-xl text-gray-400 max-w-2xl mx-auto">
          A command-line tool to create C++ projects from <code className="text-cyan-400 bg-cyan-400/10 px-2 py-0.5 rounded">forge.yaml</code> files.
          <br />Like Cargo for Rust, but for C++!
        </p>

        <div className="flex items-center justify-center gap-4">
          <a 
            href="https://github.com/ozacod/forge" 
            target="_blank" 
            rel="noopener noreferrer"
            className="flex items-center gap-2 px-6 py-3 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors text-gray-300"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
            </svg>
            View on GitHub
          </a>
        </div>
      </div>

      {/* One-liner Install Section */}
      <div className="card-glass rounded-3xl p-8 animate-fade-in border-2 border-cyan-500/30" style={{ animationDelay: '100ms' }}>
        <div className="text-center mb-6">
          <h2 className="font-display text-2xl font-bold text-white mb-2">
            ⚡ Quick Install
          </h2>
          <p className="text-gray-400">
            Auto-detects your OS and architecture. Works on macOS, Linux, and WSL.
          </p>
        </div>

        {/* Install method tabs */}
        <div className="flex items-center justify-center gap-2 mb-6">
          <button
            onClick={() => setInstallMethod('script')}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              installMethod === 'script'
                ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            Install via curl
          </button>
          <button
            onClick={() => setInstallMethod('manual')}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              installMethod === 'manual'
                ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            Install via wget
          </button>
        </div>

        {/* Install command */}
        <div className="relative max-w-4xl mx-auto">
          <pre className="code-preview p-4 rounded-xl text-sm md:text-base overflow-x-auto bg-black/50 border border-white/10">
            <code className="text-green-400 font-mono">
              {installMethod === 'script' ? INSTALL_SCRIPT_CURL : INSTALL_SCRIPT_WGET}
            </code>
          </pre>
          <button
            onClick={() => copyCommand(installMethod === 'script' ? INSTALL_SCRIPT_CURL : INSTALL_SCRIPT_WGET)}
            className="absolute top-3 right-3 p-2 text-gray-400 hover:text-white bg-black/70 rounded-lg transition-colors"
            title="Copy to clipboard"
          >
            {copiedCommand === (installMethod === 'script' ? INSTALL_SCRIPT_CURL : INSTALL_SCRIPT_WGET) ? (
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

        <p className="text-center text-gray-500 text-sm mt-4">
          The script will download the right binary for your system and install it to <code className="text-gray-400">/usr/local/bin</code>
        </p>
      </div>

      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {FEATURES.map((feature, index) => (
          <div 
            key={feature.title}
            className="card-glass rounded-2xl p-5 animate-slide-up"
            style={{ animationDelay: `${index * 100}ms` }}
          >
            <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-cyan-500/20 to-purple-500/20 flex items-center justify-center text-cyan-400 mb-4">
              {feature.icon}
            </div>
            <h3 className="font-display font-semibold text-white mb-2">{feature.title}</h3>
            <p className="text-sm text-gray-400">{feature.description}</p>
          </div>
        ))}
      </div>

      {/* Manual Download Section */}
      <div className="card-glass rounded-3xl p-8 animate-fade-in" style={{ animationDelay: '200ms' }}>
        <h2 className="font-display text-2xl font-bold text-white mb-2 text-center">
          Manual Download
        </h2>
        <p className="text-gray-400 text-center mb-6">
          Or download the binary directly for your platform
        </p>

        {/* Platform Selector */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          {PLATFORMS.map((platform) => (
            <button
              key={platform.id}
              onClick={() => setSelectedPlatform(selectedPlatform?.id === platform.id ? null : platform)}
              className={`p-6 rounded-2xl border-2 transition-all ${
                selectedPlatform?.id === platform.id
                  ? 'border-cyan-400 bg-cyan-400/10'
                  : 'border-white/10 bg-white/5 hover:border-white/20 hover:bg-white/10'
              }`}
            >
              <div className="flex flex-col items-center gap-3">
                <div className={`${selectedPlatform?.id === platform.id ? 'text-cyan-400' : 'text-gray-400'}`}>
                  {platform.icon}
                </div>
                <span className={`font-display font-semibold ${selectedPlatform?.id === platform.id ? 'text-white' : 'text-gray-300'}`}>
                  {platform.name}
                </span>
              </div>
            </button>
          ))}
        </div>

        {/* Architecture Selection */}
        {selectedPlatform && (
          <div className="space-y-4 animate-fade-in">
            <h3 className="text-lg font-semibold text-white">Select Architecture</h3>
            <div className="space-y-3">
              {selectedPlatform.architectures.map((arch) => (
                <div key={arch.id} className="bg-black/30 rounded-xl p-4 border border-white/5">
                  <div className="flex items-center justify-between mb-3">
                    <div>
                      <span className="font-mono text-cyan-400">{arch.filename}</span>
                      <span className="text-gray-500 text-sm ml-3">{arch.name}</span>
                    </div>
                    <a
                      href={`https://github.com/ozacod/forge/releases/latest/download/${arch.filename}`}
                      className="btn-primary px-4 py-2 rounded-lg font-semibold text-sm flex items-center gap-2"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                      Download
                    </a>
                  </div>
                  
                  {/* Install Command */}
                  <div className="relative">
                    <pre className="code-preview p-3 rounded-lg text-sm overflow-x-auto">
                      <code className="text-gray-300">{getInstallCommand(arch)}</code>
                    </pre>
                    <button
                      onClick={() => copyCommand(getInstallCommand(arch))}
                      className="absolute top-2 right-2 p-2 text-gray-500 hover:text-white bg-black/50 rounded-lg transition-colors"
                      title="Copy to clipboard"
                    >
                      {copiedCommand === getInstallCommand(arch) ? (
                        <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                      ) : (
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                        </svg>
                      )}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Quick Start */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Usage */}
        <div className="card-glass rounded-2xl p-6 animate-slide-up" style={{ animationDelay: '300ms' }}>
          <h3 className="font-display text-xl font-bold text-white mb-4 flex items-center gap-2">
            <svg className="w-5 h-5 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            Quick Start
          </h3>
          <div className="space-y-4">
            <div className="space-y-2">
              <span className="text-xs font-mono text-gray-500">Create & run project</span>
              <CodeBlock code="forge new my_app\ncd my_app\nforge generate\nforge run" />
            </div>
            <div className="space-y-2">
              <span className="text-xs font-mono text-gray-500">Add dependencies</span>
              <CodeBlock code="forge add spdlog\nforge generate" />
            </div>
            <div className="space-y-2">
              <span className="text-xs font-mono text-gray-500">Build & test</span>
              <CodeBlock code="forge build --release\nforge test" />
            </div>
          </div>
        </div>

        {/* forge.yaml Example */}
        <div className="card-glass rounded-2xl p-6 animate-slide-up" style={{ animationDelay: '400ms' }}>
          <h3 className="font-display text-xl font-bold text-white mb-4 flex items-center gap-2">
            <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            forge.yaml
          </h3>
          <pre className="code-preview p-4 rounded-xl text-sm overflow-x-auto">
            <code className="text-gray-300">{`package:
  name: my_project
  version: "0.1.0"
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

dependencies:
  spdlog:
    spdlog_header_only: true
  nlohmann_json: {}
  fmt: {}

dev-dependencies:
  catch2: {}`}</code>
          </pre>
        </div>
      </div>

      {/* Available Commands */}
      <div className="card-glass rounded-2xl p-6 animate-fade-in" style={{ animationDelay: '500ms' }}>
        <h3 className="font-display text-xl font-bold text-white mb-6 flex items-center gap-2">
          <svg className="w-5 h-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Available Commands
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Project Commands */}
          <div>
            <h4 className="text-sm font-semibold text-cyan-400 mb-3 uppercase tracking-wider">Project</h4>
            <div className="space-y-2 text-sm">
              <div><code className="text-green-400">new &lt;name&gt;</code> <span className="text-gray-400">- Create new project</span></div>
              <div><code className="text-green-400">new &lt;name&gt; --lib</code> <span className="text-gray-400">- Create library project</span></div>
              <div><code className="text-green-400">init</code> <span className="text-gray-400">- Init in current directory</span></div>
              <div><code className="text-green-400">init -t &lt;template&gt;</code> <span className="text-gray-400">- Use template</span></div>
            </div>
          </div>

          {/* Build Commands */}
          <div>
            <h4 className="text-sm font-semibold text-cyan-400 mb-3 uppercase tracking-wider">Generate & Build</h4>
            <div className="space-y-2 text-sm">
              <div><code className="text-green-400">generate</code> <span className="text-gray-400">- Generate CMake from yaml</span></div>
              <div><code className="text-green-400">build</code> <span className="text-gray-400">- Compile (Debug mode)</span></div>
              <div><code className="text-green-400">build --release</code> <span className="text-gray-400">- Release mode (O2)</span></div>
              <div><code className="text-green-400">build -O3</code> <span className="text-gray-400">- Optimize with O3</span></div>
              <div><code className="text-green-400">build --clean</code> <span className="text-gray-400">- Clean and rebuild</span></div>
              <div><code className="text-green-400">run</code> <span className="text-gray-400">- Build and run executable</span></div>
              <div><code className="text-green-400">test</code> <span className="text-gray-400">- Build and run tests</span></div>
              <div><code className="text-green-400">clean</code> <span className="text-gray-400">- Remove build artifacts</span></div>
            </div>
          </div>

          {/* Dependency Commands */}
          <div>
            <h4 className="text-sm font-semibold text-cyan-400 mb-3 uppercase tracking-wider">Dependencies</h4>
            <div className="space-y-2 text-sm">
              <div><code className="text-green-400">add &lt;lib&gt;</code> <span className="text-gray-400">- Add dependency</span></div>
              <div><code className="text-green-400">add --dev &lt;lib&gt;</code> <span className="text-gray-400">- Add dev dependency</span></div>
              <div><code className="text-green-400">remove &lt;lib&gt;</code> <span className="text-gray-400">- Remove dependency</span></div>
              <div><code className="text-green-400">update</code> <span className="text-gray-400">- Update all dependencies</span></div>
              <div><code className="text-green-400">list</code> <span className="text-gray-400">- List available libraries</span></div>
              <div><code className="text-green-400">search &lt;query&gt;</code> <span className="text-gray-400">- Search libraries</span></div>
              <div><code className="text-green-400">info &lt;lib&gt;</code> <span className="text-gray-400">- Show library details</span></div>
            </div>
          </div>

          {/* Development Commands */}
          <div>
            <h4 className="text-sm font-semibold text-cyan-400 mb-3 uppercase tracking-wider">Development</h4>
            <div className="space-y-2 text-sm">
              <div><code className="text-green-400">fmt</code> <span className="text-gray-400">- Format code (clang-format)</span></div>
              <div><code className="text-green-400">fmt --check</code> <span className="text-gray-400">- Check formatting</span></div>
              <div><code className="text-green-400">lint</code> <span className="text-gray-400">- Run static analysis (clang-tidy)</span></div>
              <div><code className="text-green-400">lint --fix</code> <span className="text-gray-400">- Auto-fix lint issues</span></div>
              <div><code className="text-green-400">doc</code> <span className="text-gray-400">- Generate documentation</span></div>
              <div><code className="text-green-400">doc --open</code> <span className="text-gray-400">- Open docs in browser</span></div>
              <div><code className="text-green-400">release &lt;patch|minor|major&gt;</code> <span className="text-gray-400">- Bump version</span></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Helper component for code blocks
function CodeBlock({ code }: { code: string }) {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative">
      <pre className="code-preview p-3 rounded-lg text-sm overflow-x-auto">
        <code className="text-gray-300">{code}</code>
      </pre>
      <button
        onClick={copy}
        className="absolute top-2 right-2 p-1.5 text-gray-500 hover:text-white bg-black/50 rounded transition-colors"
      >
        {copied ? (
          <svg className="w-3.5 h-3.5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        ) : (
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        )}
      </button>
    </div>
  );
}

