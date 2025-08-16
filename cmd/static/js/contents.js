const ICONS = {
    'twitter': 'twitter', 'youtube': 'youtube', 'github': 'github', 'linkedin': 'linkedin',
    'folder': 'folder', 'dot': 'dot', 'figma': 'figma', 'code': 'code', 'package': 'package',
    'terminal': 'terminal-square', 'newspaper': 'newspaper', 'pen': 'pen-square', 'rss': 'rss',
    'default': 'link-2'
};

const getIcon = (name) => ICONS[name] || ICONS['default'];

const createLinkHTML = (link) => `
    <a href="${link.url}" class="flex items-center text-ctp-subtext1 hover:text-ctp-rosewater transition-colors">
        <i data-lucide="${getIcon(link.icon)}" class="w-4 h-4 mr-3"></i>${link.name}
    </a>`;

const createFolderHTML = (folder) => `
    <details open>
        <summary class="flex items-center text-ctp-subtext1 hover:text-ctp-rosewater transition-colors mb-2">
            <i data-lucide="chevron-right" class="chevron-icon w-4 h-4 mr-2"></i>
            <i data-lucide="${getIcon(folder.icon)}" class="w-4 h-4 mr-2 text-ctp-peach"></i>
            ${folder.name}
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
    <img src="${header.logoURL}" alt="Logo" class="w-16 h-16 rounded-full mr-4" onerror="this.onerror=null;this.src='[https://placehold.co/40x40/181825/cdd6f4?text=](https://placehold.co/40x40/181825/cdd6f4?text=):)';">
    <h1 class="text-3xl font-bold text-ctp-text">${header.title}</h1>`;

const createWeatherHTML = (weather) => weather ? `
    <i data-lucide="sun" class="w-6 h-6 text-ctp-yellow"></i>
    <span class="font-semibold text-lg">${weather.tempC}°C</span>
    <span class="text-sm">${weather.description}</span>` : '';

// Card Definitions
const CARD_TEMPLATES = {
    github: (repo) => `
        <div class="bg-ctp-base p-4 rounded-lg">
            <a href="${repo.url}" class="font-semibold text-ctp-blue hover:underline mb-1 block truncate">${repo.name}</a>
            <div class="flex items-center space-x-4 text-sm text-ctp-subtext1 mt-2">
                <span class="flex items-center"><i data-lucide="star" class="w-4 h-4 mr-1.5 text-ctp-yellow"></i>${repo.stars}</span>
                <span class="flex items-center"><i data-lucide="git-fork" class="w-4 h-4 mr-1.5 text-ctp-teal"></i>${repo.forks}</span>
            </div>
        </div>`,
    
    youtube: (video) => `
        <a href="${video.url}" class="group block">
            <img src="${video.thumbnail}" alt="Video Thumbnail" class="w-full h-28 object-cover rounded-lg mb-2" onerror="this.onerror=null;this.src='[https://placehold.co/600x400/181825/cdd6f4?text=Video](https://placehold.co/600x400/181825/cdd6f4?text=Video)';">
            <h3 class="font-semibold text-sm text-ctp-text truncate group-hover:text-ctp-red transition-colors">${video.title}</h3>
            <p class="text-xs text-ctp-subtext0">${video.channel}</p>
        </a>`,

    reddit: (post) => `<!-- Reddit card placeholder -->`,
    service_stats: (service) => `<!-- Service stats card placeholder -->`,
};

// Section Definitions
const createStatusCardSectionHTML = (section) => `
    <section>
        <div class="h-px bg-ctp-surface0 my-8 w-1/4 mx-auto"></div>
        <h2 class="text-xl font-bold text-ctp-lavender mb-4 flex items-center">
            <i data-lucide="${getIcon(section.icon)}" class="mr-3"></i>${section.title}
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            ${section.cards.map(CARD_TEMPLATES.github).join('')}
        </div>
    </section>`;

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
