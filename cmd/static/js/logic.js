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

const setupSettingsModal = () => {
    const modal = document.getElementById('settings-modal');
    const closeBtn = document.getElementById('close-modal-btn');
    const cancelBtn = document.getElementById('cancel-btn');
    const saveBtn = document.getElementById('save-btn');
    const editor = document.getElementById('config-editor');
    const saveStatus = document.getElementById('save-status');
    const saveBtnText = document.getElementById('save-btn-text');
    const saveSpinner = document.getElementById('save-spinner');

    const openModal = async () => {
        saveStatus.textContent = '';
        editor.value = 'Loading...';
        modal.classList.remove('hidden');
        modal.classList.add('flex');
        try {
            const response = await fetch('/api/config');
            if (!response.ok) throw new Error('Failed to fetch config');
            const configText = await response.text();
            editor.value = configText;
        } catch (error) {
            console.error('Error opening settings:', error);
            editor.value = `# Error loading config:\n${error.message}`;
        }
    };

    const closeModal = () => {
        modal.classList.add('hidden');
        modal.classList.remove('flex');
    };

    const saveConfig = async () => {
        saveBtn.disabled = true;
        saveBtnText.textContent = 'Saving...';
        saveSpinner.classList.remove('hidden');
        saveStatus.textContent = '';
        try {
            const response = await fetch('/api/config', {
                method: 'POST',
                headers: { 'Content-Type': 'application/yaml' },
                body: editor.value
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to save config');
            }
            saveStatus.textContent = 'Configuration saved successfully! Reloading dashboard...';
            setTimeout(async () => {
                closeModal();
                await Promise.all([
                    renderHeader(),
                    renderBookmarks(),
                    renderDynamicSections(),
                ]);
                lucide.createIcons();
            }, 1500);
        } catch (error) {
            console.error('Error saving config:', error);
            saveStatus.textContent = `Error: ${error.message}`;
        } finally {
            saveBtn.disabled = false;
            saveBtnText.textContent = 'Save';
            saveSpinner.classList.add('hidden');
        }
    };

    document.body.addEventListener('click', (e) => {
        if (e.target.closest('#settings-btn')) {
            openModal();
        }
    });
    closeBtn.addEventListener('click', closeModal);
    cancelBtn.addEventListener('click', closeModal);
    saveBtn.addEventListener('click', saveConfig);
};

document.addEventListener('DOMContentLoaded', () => {
    setupSidebar();
    setupSettingsModal();
    Promise.all([
        renderHeader(),
        renderBookmarks(),
        renderDynamicSections(),
    ]).then(() => {
        lucide.createIcons();
    });
});
