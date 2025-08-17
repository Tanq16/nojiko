const ICONS = {
    'twitter': 'twitter', 'youtube': 'youtube', 'github': 'github', 'linkedin': 'linkedin',
    'folder': 'folder', 'dot': 'dot', 'figma': 'figma', 'code': 'code', 'package': 'package',
    'terminal-square': 'terminal-square', 'newspaper': 'newspaper', 'pen-square': 'pen-square', 'rss': 'rss',
    'bar-chart-3': 'bar-chart-3', 'tv-minimal': 'tv-minimal',
    'default': 'link-2'
};

const getIcon = (name) => ICONS[name] || ICONS['default'];

const createLinkHTML = (link) => `
    <a href="${link.url}" class="flex items-center text-ctp-subtext1 hover:text-ctp-rosewater transition-colors">
        <i data-lucide="${getIcon(link.icon)}" class="w-4 h-4 mr-3"></i>${link.name || 'Unnamed'}
    </a>`;

const createFolderHTML = (folder) => `
    <details open>
        <summary class="flex items-center text-ctp-subtext1 hover:text-ctp-rosewater transition-colors mb-2">
            <i data-lucide="chevron-right" class="chevron-icon w-4 h-4 mr-2"></i>
            <i data-lucide="${getIcon(folder.icon)}" class="w-4 h-4 mr-2 text-ctp-peach"></i>
            ${folder.name || 'Unnamed'}
        </summary>
        <div class="pl-6 space-y-2">
            ${folder.links.map(createLinkHTML).join('')}
        </div>
    </details>`;

const createCategoryHTML = (cat) => `
    <details open>
        <summary class="font-semibold text-${cat.color || 'ctp-green'} mb-2 flex items-center">
            <i data-lucide="chevron-right" class="chevron-icon w-4 h-4 mr-2"></i>${cat.category}
        </summary>
        <div class="space-y-2 pl-6">
            ${(cat.links || []).map(createLinkHTML).join('')}
            ${(cat.folders || []).map(createFolderHTML).join('')}
        </div>
    </details>`;

const createHeaderContentHTML = (header) => `
    ${header.showLogo ? `<img src="${header.logoURL}" alt="Logo" class="w-16 h-16 rounded-full mr-4" onerror="this.onerror=null;this.src='logo-low.png';">` : ''}
    <h1 class="text-3xl font-bold text-ctp-text">${header.title}</h1>`;

const createWeatherHTML = (weather) => weather ? `
    <a href="https://open-meteo.com/" target="_blank" rel="noopener noreferrer" class="flex flex-col items-center sm:items-end w-full">
        <div class="flex items-center space-x-3 text-ctp-subtext1">
            <i data-lucide="sun" class="w-6 h-6 text-ctp-yellow"></i>
            <span class="font-semibold text-lg">${weather.tempC}°C</span>
            <span class="text-sm">${weather.description}</span>
        </div>
        <p class="text-xs text-ctp-overlay0">via Open-Meteo.com</p>
    </a>` : '';

const CARD_TEMPLATES = {
    github: (repo) => `
        <div class="bg-ctp-base p-4 rounded-lg">
            <a href="${repo.url}" target="_blank" rel="noopener noreferrer" class="font-semibold text-ctp-blue hover:underline mb-1 block truncate">${repo.name}</a>
            <div class="flex items-center space-x-4 text-sm text-ctp-subtext1 mt-2">
                <span class="flex items-center" title="Stars"><i data-lucide="star" class="w-4 h-4 mr-1.5 text-ctp-yellow"></i>${repo.stars}</span>
                <span class="flex items-center" title="Issues"><i data-lucide="alert-circle" class="w-4 h-4 mr-1.5 text-ctp-red"></i>${repo.issues}</span>
                <span class="flex items-center" title="Pull Requests"><i data-lucide="git-pull-request" class="w-4 h-4 mr-1.5 text-ctp-green"></i>${repo.prs}</span>
            </div>
        </div>`,
    service: (service) => `
        <div class="bg-ctp-base p-4 rounded-lg flex items-center justify-between">
            <span class="font-semibold text-ctp-text">${service.name}</span>
            <div class="w-3 h-3 rounded-full bg-ctp-green" title="Operational"></div>
        </div>`,
    youtube: (video) => `
        <a href="${video.url}" target="_blank" rel="noopener noreferrer" class="group block bg-ctp-base rounded-lg overflow-hidden">
            <div class="relative">
                <img src="${video.thumbnail}" alt="Video Thumbnail" class="w-full h-32 object-cover" onerror="this.onerror=null;this.src='https://placehold.co/600x400/181825/cdd6f4?text=Video';">
                <div class="absolute inset-0 bg-black/20 group-hover:bg-black/40 transition-all"></div>
            </div>
            <div class="p-3">
                <h3 class="font-semibold text-sm text-ctp-text truncate group-hover:text-ctp-red transition-colors" title="${video.title}">${video.title}</h3>
                <div class="flex justify-between items-center mt-1">
                    <p class="text-xs text-ctp-subtext0 truncate">@${video.channel}</p>
                    <p class="text-xs text-ctp-overlay2 flex-shrink-0 ml-2">${video.published}</p>
                </div>
            </div>
        </a>`,
};

const createStatusCardSectionHTML = (section) => {
    const cardTemplate = CARD_TEMPLATES[section.type] || (() => '<div>Unsupported card type</div>');
    return `
    <section>
        <div class="h-px bg-ctp-surface0 my-8 w-1/4 mx-auto"></div>
        <h2 class="text-xl font-bold text-ctp-lavender mb-4 flex items-center">
            <i data-lucide="${getIcon(section.icon)}" class="mr-3"></i>${section.title}
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            ${section.cards.map(cardTemplate).join('')}
        </div>
    </section>`;
}

const createThumbFeedSectionHTML = (section) => {
    const cardTemplate = CARD_TEMPLATES[section.feedType] || (() => '<div>Unsupported feed type</div>');
    return `
    <section>
        <div class="h-px bg-ctp-surface0 my-8 w-1/4 mx-auto"></div>
        <h2 class="text-xl font-bold text-ctp-red mb-4 flex items-center">
            <i data-lucide="${getIcon(section.icon)}" class="mr-3"></i>${section.title}
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
            ${section.cards.map(cardTemplate).join('')}
        </div>
    </section>`;
}
