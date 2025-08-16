const fetchData = async (url) => {
    try {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        return await response.json();
    } catch (error) {
        console.error(`Failed to fetch from ${url}:`, error);
        return null;
    }
};

const renderBookmarks = async () => {
    const data = await fetchData('/api/bookmarks');
    if (data) {
        document.getElementById('bookmarks-container').innerHTML = data.map(createCategoryHTML).join('');
    }
};

const renderHeader = async () => {
    const data = await fetchData('/api/header');
    if (data) {
        document.getElementById('header-content').innerHTML = createHeaderContentHTML(data);
        document.getElementById('weather-content').innerHTML = createWeatherHTML(data.weather);
    }
};

const renderDynamicSections = async () => {
    const container = document.getElementById('dynamic-sections-container');
    container.innerHTML = '';
    const statusCards = await fetchData('/api/status-cards');
    if (statusCards) {
        container.innerHTML += statusCards.map(createStatusCardSectionHTML).join('');
    }
    const thumbFeeds = await fetchData('/api/thumb-feeds');
    if (thumbFeeds) {
        container.innerHTML += thumbFeeds.map(createThumbFeedSectionHTML).join('');
    }
};

const setupSidebar = () => {
    const sidebar = document.getElementById('sidebar');
    const overlay = document.getElementById('sidebar-overlay');
    const openBtn = document.getElementById('open-sidebar-btn');
    const closeBtn = document.getElementById('close-sidebar-btn');
    const openSidebar = () => {
        sidebar.classList.remove('-translate-x-full');
        overlay.classList.remove('hidden');
    };
    const closeSidebar = () => {
        sidebar.classList.add('-translate-x-full');
        overlay.classList.add('hidden');
    };
    openBtn.addEventListener('click', openSidebar);
    closeBtn.addEventListener('click', closeSidebar);
    overlay.addEventListener('click', closeSidebar);
};

document.addEventListener('DOMContentLoaded', () => {
    setupSidebar();
    Promise.all([
        renderHeader(),
        renderBookmarks(),
        renderDynamicSections(),
    ]).then(() => {
        lucide.createIcons();
    });
});
