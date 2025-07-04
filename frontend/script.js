// JavaScript for frontend logic
// Future steps: Add AJAX to fetch /api/fights, implement search and filtering

// Global variables
let currentFights = [];
let currentPage = 1;
let fightsPerPage = 10;

// DOM elements
const fightsBody = document.getElementById('fightsBody');
const loadingDiv = document.getElementById('loading');
const refreshBtn = document.getElementById('refreshBtn');
const searchInput = document.getElementById('searchInput');
const filterSelect = document.getElementById('filterSelect');

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    console.log('EasyPars frontend loaded');
    
    // Set up event listeners
    setupEventListeners();
    
    // Load initial data
    loadFights();
});

// Set up event listeners
function setupEventListeners() {
    // Refresh button
    refreshBtn.addEventListener('click', loadFights);
    
    // Search input
    // Future steps: Add debounced search functionality
    searchInput.addEventListener('input', function() {
        console.log('Search functionality - to be implemented');
    });
    
    // Filter select
    // Future steps: Add filtering by result type
    filterSelect.addEventListener('change', function() {
        console.log('Filter functionality - to be implemented');
    });
}

// Load fights from API
async function loadFights() {
    try {
        showLoading(true);
        
        // Fetch data from API
        const response = await fetch('/api/fights');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        currentFights = data.data || [];
        
        // Display fights
        displayFights(currentFights);
        
        console.log('Fights loaded successfully:', currentFights.length);
        
    } catch (error) {
        console.error('Error loading fights:', error);
        displayError('Failed to load fights. Please try again.');
    } finally {
        showLoading(false);
    }
}

// Display fights in the table
function displayFights(fights) {
    if (!fights || fights.length === 0) {
        fightsBody.innerHTML = '<tr><td colspan="7">No fights found</td></tr>';
        return;
    }
    
    const html = fights.map(fight => `
        <tr>
            <td>${fight.date}</td>
            <td>${fight.fighter1}</td>
            <td>${fight.fighter2}</td>
            <td>${fight.result}</td>
            <td>${fight.location}</td>
            <td>${fight.round || 'N/A'}</td>
            <td>${fight.time || 'N/A'}</td>
        </tr>
    `).join('');
    
    fightsBody.innerHTML = html;
}

// Show/hide loading indicator
function showLoading(show) {
    loadingDiv.style.display = show ? 'block' : 'none';
}

// Display error message
function displayError(message) {
    fightsBody.innerHTML = `<tr><td colspan="7" class="error">${message}</td></tr>`;
}

// Future functions to be implemented:
// - searchFights(query)
// - filterFights(filter)
// - paginateFights(page)
// - sortFights(column, direction)
// - exportFights(format)
// - showFightDetails(fightId)