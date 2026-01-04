// ===========================
// Global State
// ===========================
let searchIndex = [];
let fuse = null;
let chatOpen = true;
let contentData = null;

// ===========================
// Initialization
// ===========================
document.addEventListener('DOMContentLoaded', () => {
    loadContent();
    initializeSidebar();
    initializeSearch();
    initializeChat();
    initializeNavigation();
});

// ===========================
// Content Loading
// ===========================
async function loadContent() {
    try {
        const response = await fetch('/docs/data/content.json');
        contentData = await response.json();
        
        renderNavigation(contentData.sections);
        renderContent(contentData.sections);
        highlightActiveSection();
        
        console.log(`Loaded ${contentData.metadata.total_sections} sections`);
    } catch (error) {
        console.error('Failed to load content:', error);
        document.getElementById('documentationContent').innerHTML = 
            '<div class="error-message">Failed to load documentation. Please try again later.</div>';
    }
}

function renderNavigation(sections) {
    const navMenu = document.getElementById('navMenu');
    navMenu.innerHTML = sections.map(section => `
        <li class="nav-item nav-level-${section.level}">
            <a href="#${section.id}" class="nav-link">${escapeHtml(section.heading)}</a>
        </li>
    `).join('');
}

function renderContent(sections) {
    const contentContainer = document.getElementById('documentationContent');
    contentContainer.innerHTML = sections.map(section => `
        <section class="doc-section" id="${section.id}">
            <h${section.level} class="section-heading">${escapeHtml(section.heading)}</h${section.level}>
            
            <div class="section-content">
                ${formatContentHTML(section.content)}
            </div>

            ${section.images && section.images.length > 0 ? `
                <div class="section-images">
                    ${section.images.map(img => `
                        <div class="image-container">
                            <img src="/docs/images/${img}" alt="${img}" loading="lazy">
                        </div>
                    `).join('')}
                </div>
            ` : ''}
        </section>
    `).join('');
}

function formatContentHTML(content) {
    if (!content) return '';
    
    // Convert markdown to HTML
    let html = content;
    
    // Handle code blocks first (```...```)
    html = html.replace(/```(\w+)?\n([\s\S]*?)```/g, (match, lang, code) => {
        const language = lang || '';
        return `<pre><code class="language-${language}">${escapeHtml(code.trim())}</code></pre>`;
    });
    
    // Handle inline code (`...`)
    html = html.replace(/`([^`]+)`/g, '<code>$1</code>');
    
    // Handle bold (**text** or __text__)
    html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>');
    
    // Handle italic (*text* or _text_)
    html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>');
    html = html.replace(/_([^_]+)_/g, '<em>$1</em>');
    
    // Handle links [text](url)
    html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');
    
    // Handle bullet points (lines starting with -, *, or +)
    const lines = html.split('\n');
    let inList = false;
    let processedLines = [];
    
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim();
        
        // Check if line is a list item
        if (/^[-*+]\s+/.test(line)) {
            if (!inList) {
                processedLines.push('<ul>');
                inList = true;
            }
            const itemContent = line.replace(/^[-*+]\s+/, '');
            processedLines.push(`<li>${itemContent}</li>`);
        } else if (line.match(/^\d+\.\s+/)) {
            // Numbered list
            if (!inList) {
                processedLines.push('<ol>');
                inList = true;
            }
            const itemContent = line.replace(/^\d+\.\s+/, '');
            processedLines.push(`<li>${itemContent}</li>`);
        } else {
            if (inList) {
                processedLines.push('</ul>');
                inList = false;
            }
            if (line) {
                processedLines.push(`<p>${line}</p>`);
            }
        }
    }
    
    if (inList) {
        processedLines.push('</ul>');
    }
    
    return processedLines.join('\n');
}

// ===========================
// Sidebar Navigation
// ===========================
function initializeSidebar() {
    const sidebar = document.getElementById('sidebar');
    const sidebarToggle = document.getElementById('sidebarToggle');
    const mainContent = document.getElementById('mainContent');

    if (sidebarToggle) {
        sidebarToggle.addEventListener('click', () => {
            sidebar.classList.toggle('active');
        });
    }

    // Close sidebar when clicking outside on mobile
    mainContent.addEventListener('click', () => {
        if (window.innerWidth <= 768 && sidebar.classList.contains('active')) {
            sidebar.classList.remove('active');
        }
    });
}

// ===========================
// Search Functionality
// ===========================
async function initializeSearch() {
    try {
        // Try to load from new structured data format first
        let response = await fetch('/docs/data/content.json');
        if (response.ok) {
            const contentData = await response.json();
            // Convert to search index format
            searchIndex = contentData.sections.map(section => ({
                id: section.id,
                heading: section.heading,
                content: section.content,
                level: section.level
            }));
        } else {
            // Fallback to old search-index.json
            response = await fetch('/docs/search-index.json');
            const data = await response.json();
            searchIndex = data.items;
        }

        // Initialize Fuse.js
        const options = {
            keys: ['heading', 'content'],
            threshold: 0.4,
            includeScore: true,
            minMatchCharLength: 2
        };
        fuse = new Fuse(searchIndex, options);

        // Setup search input
        const searchInput = document.getElementById('searchInput');
        const searchResults = document.getElementById('searchResults');

        searchInput.addEventListener('input', (e) => {
            const query = e.target.value.trim();

            if (query.length < 2) {
                searchResults.classList.remove('active');
                searchResults.innerHTML = '';
                return;
            }

            performSearch(query);
        });

        // Close search results when clicking outside
        document.addEventListener('click', (e) => {
            if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
                searchResults.classList.remove('active');
            }
        });

    } catch (error) {
        console.error('Failed to initialize search:', error);
    }
}

function performSearch(query) {
    const searchResults = document.getElementById('searchResults');
    const results = fuse.search(query);

    if (results.length === 0) {
        searchResults.innerHTML = '<div class="search-result-item"><div class="search-result-heading">No results found</div></div>';
        searchResults.classList.add('active');
        return;
    }

    // Display results
    const html = results.slice(0, 5).map(result => {
        const item = result.item;
        return `
            <div class="search-result-item" onclick="navigateToSection('${item.id}')">
                <div class="search-result-heading">${escapeHtml(item.heading)}</div>
                <div class="search-result-content">${escapeHtml(item.content)}</div>
            </div>
        `;
    }).join('');

    searchResults.innerHTML = html;
    searchResults.classList.add('active');
}

function navigateToSection(sectionId) {
    const section = document.getElementById(sectionId);
    if (section) {
        section.scrollIntoView({ behavior: 'smooth', block: 'start' });
        
        // Close search results
        const searchResults = document.getElementById('searchResults');
        searchResults.classList.remove('active');
        
        // Clear search input
        document.getElementById('searchInput').value = '';

        // Close sidebar on mobile
        if (window.innerWidth <= 768) {
            document.getElementById('sidebar').classList.remove('active');
        }
    }
}

// ===========================
// Chat Widget
// ===========================
function initializeChat() {
    const chatWidget = document.getElementById('chatWidget');
    const chatFab = document.getElementById('chatFab');
    const chatToggle = document.getElementById('chatToggle');
    const chatForm = document.getElementById('chatForm');

    // Toggle chat
    chatToggle.addEventListener('click', () => {
        chatOpen = false;
        chatWidget.classList.add('minimized');
        chatFab.classList.add('active');
    });

    chatFab.addEventListener('click', () => {
        chatOpen = true;
        chatWidget.classList.remove('minimized');
        chatFab.classList.remove('active');
    });

    // Handle chat form submission
    chatForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await handleChatSubmit();
    });
}

async function handleChatSubmit() {
    const chatInput = document.getElementById('chatInput');
    const chatMessages = document.getElementById('chatMessages');
    const prompt = chatInput.value.trim();

    if (!prompt) return;

    // Add user message
    addChatMessage(prompt, 'user');
    chatInput.value = '';

    // Show typing indicator
    const typingIndicator = document.createElement('div');
    typingIndicator.className = 'chat-message bot-message typing-indicator';
    typingIndicator.innerHTML = '<span></span><span></span><span></span>';
    typingIndicator.id = 'typingIndicator';
    chatMessages.appendChild(typingIndicator);
    scrollChatToBottom();

    try {
        // Send request to backend
        const response = await fetch('/api/chat', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ prompt })
        });

        const data = await response.json();

        // Remove typing indicator
        typingIndicator.remove();

        if (data.error) {
            addChatMessage(`Error: ${data.error}`, 'bot');
        } else {
            addChatMessage(data.answer, 'bot');
        }

    } catch (error) {
        // Remove typing indicator
        typingIndicator.remove();
        addChatMessage('Failed to get response. Make sure Ollama is running.', 'bot');
        console.error('Chat error:', error);
    }
}

function addChatMessage(text, type) {
    const chatMessages = document.getElementById('chatMessages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `chat-message ${type}-message`;

    // Format the message with proper HTML
    const formattedText = formatChatMessage(text);
    messageDiv.innerHTML = formattedText;

    chatMessages.appendChild(messageDiv);
    scrollChatToBottom();
}

function formatChatMessage(text) {
    if (!text) return '';
    
    let formatted = text;
    
    // First, handle code blocks (triple backticks) - escape content
    formatted = formatted.replace(/```([\w]*)?[\n\r]?([\s\S]*?)```/g, (match, lang, code) => {
        const escapedCode = code
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;');
        return `<pre><code class="language-${lang || ''}">${escapedCode.trim()}</code></pre>`;
    });
    
    // Handle inline code (single backticks) - escape content
    formatted = formatted.replace(/`([^`]+)`/g, (match, code) => {
        const escapedCode = code
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;');
        return `<code>${escapedCode}</code>`;
    });
    
    // Escape remaining HTML (but not our formatted tags)
    formatted = formatted.split(/(<pre>[\s\S]*?<\/pre>|<code>[\s\S]*?<\/code>)/).map((part, i) => {
        if (i % 2 === 0) {
            return part.replace(/&/g, '&amp;')
                      .replace(/</g, '&lt;')
                      .replace(/>/g, '&gt;');
        }
        return part;
    }).join('');
    
    // Convert **bold** text
    formatted = formatted.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    
    // Convert *italic* text (single asterisk, not double)
    formatted = formatted.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>');

    // Split into lines for list processing
    const lines = formatted.split('\n');
    const result = [];
    let inOrderedList = false;
    let inUnorderedList = false;
    let currentParagraph = [];

    for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim();

        // Ordered list (1. 2. 3. etc.)
        if (/^\d+\.\s/.test(line)) {
            // Close any open paragraph
            if (currentParagraph.length > 0) {
                result.push('<p>' + currentParagraph.join(' ') + '</p>');
                currentParagraph = [];
            }

            // Start ordered list if not already in one
            if (!inOrderedList) {
                if (inUnorderedList) {
                    result.push('</ul>');
                    inUnorderedList = false;
                }
                result.push('<ol>');
                inOrderedList = true;
            }

            // Add list item (remove the number prefix)
            const content = line.replace(/^\d+\.\s/, '');
            result.push('<li>' + content + '</li>');
        }
        // Unordered list (- or * at start)
        else if (/^[-*]\s/.test(line)) {
            // Close any open paragraph
            if (currentParagraph.length > 0) {
                result.push('<p>' + currentParagraph.join(' ') + '</p>');
                currentParagraph = [];
            }

            // Start unordered list if not already in one
            if (!inUnorderedList) {
                if (inOrderedList) {
                    result.push('</ol>');
                    inOrderedList = false;
                }
                result.push('<ul>');
                inUnorderedList = true;
            }

            // Add list item (remove the - or * prefix)
            const content = line.replace(/^[-*]\s/, '');
            result.push('<li>' + content + '</li>');
        }
        // Empty line - close lists and paragraphs
        else if (line === '') {
            if (inOrderedList) {
                result.push('</ol>');
                inOrderedList = false;
            }
            if (inUnorderedList) {
                result.push('</ul>');
                inUnorderedList = false;
            }
            if (currentParagraph.length > 0) {
                result.push('<p>' + currentParagraph.join(' ') + '</p>');
                currentParagraph = [];
            }
        }
        // Regular text
        else {
            // Close lists if we were in one
            if (inOrderedList) {
                result.push('</ol>');
                inOrderedList = false;
            }
            if (inUnorderedList) {
                result.push('</ul>');
                inUnorderedList = false;
            }

            // Add to current paragraph
            currentParagraph.push(line);
        }
    }

    // Close any remaining open tags
    if (inOrderedList) {
        result.push('</ol>');
    }
    if (inUnorderedList) {
        result.push('</ul>');
    }
    if (currentParagraph.length > 0) {
        result.push('<p>' + currentParagraph.join(' ') + '</p>');
    }

    return result.join('');
}

function scrollChatToBottom() {
    const chatBody = document.getElementById('chatBody');
    chatBody.scrollTop = chatBody.scrollHeight;
}

// ===========================
// Navigation Highlighting
// ===========================
function initializeNavigation() {
    const navLinks = document.querySelectorAll('.nav-link');

    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const targetId = link.getAttribute('href').substring(1);
            navigateToSection(targetId);
            
            // Update active state
            navLinks.forEach(l => l.classList.remove('active'));
            link.classList.add('active');
        });
    });

    // Highlight navigation on scroll
    window.addEventListener('scroll', highlightActiveSection);
}

function highlightActiveSection() {
    const sections = document.querySelectorAll('.doc-section');
    const navLinks = document.querySelectorAll('.nav-link');
    
    let currentSection = '';
    const scrollPosition = window.scrollY + 100;

    sections.forEach(section => {
        const sectionTop = section.offsetTop;
        const sectionHeight = section.offsetHeight;

        if (scrollPosition >= sectionTop && scrollPosition < sectionTop + sectionHeight) {
            currentSection = section.getAttribute('id');
        }
    });

    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href') === `#${currentSection}`) {
            link.classList.add('active');
        }
    });
}

// ===========================
// Utility Functions
// ===========================
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Handle keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + K to focus search
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        document.getElementById('searchInput').focus();
    }

    // Escape to close search results
    if (e.key === 'Escape') {
        const searchResults = document.getElementById('searchResults');
        searchResults.classList.remove('active');
    }
});
