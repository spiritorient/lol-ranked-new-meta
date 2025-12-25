// API URL - automatically uses same origin (works when frontend is served from backend)
const API_URL = window.location.origin;

let matchData = null;

async function analyzeMatch() {
    const region = document.getElementById('region').value.trim();
    const gameId = document.getElementById('gameId').value.trim();
    const championName = document.getElementById('championName').value.trim();
    const summonerName = document.getElementById('summonerName').value.trim();

    if (!region) {
        showError('Please select a Region');
        return;
    }

    if (!gameId) {
        showError('Please enter a Game ID');
        return;
    }

    // Combine region and game ID to create full match ID
    const matchId = `${region}_${gameId}`;

    // Show loading
    document.getElementById('loading').classList.remove('hidden');
    document.getElementById('error').classList.add('hidden');
    document.getElementById('results').classList.add('hidden');

    try {
        // Build URL with query params
        let url = `${API_URL}/analyze-match-get?match_id=${encodeURIComponent(matchId)}`;
        if (championName) {
            url += `&champion_name=${encodeURIComponent(championName)}`;
        }
        if (summonerName) {
            url += `&summoner_name=${encodeURIComponent(summonerName)}`;
        }

        const response = await fetch(url);
        const data = await response.json();

        if (data.error) {
            showError(data.error);
            return;
        }

        matchData = data;
        displayResults(data);
    } catch (error) {
        showError(`Failed to analyze match: ${error.message}`);
    } finally {
        document.getElementById('loading').classList.add('hidden');
    }
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
}

function displayResults(data) {
    document.getElementById('results').classList.remove('hidden');
    switchTab('overview');
}

function switchTab(tabName) {
    // Update tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

    // Update content
    const contentDiv = document.getElementById('tab-content');
    
    switch(tabName) {
        case 'overview':
            contentDiv.innerHTML = renderOverview(matchData);
            break;
        case 'what-went-well':
            contentDiv.innerHTML = renderWhatWentWell(matchData);
            break;
        case 'what-went-wrong':
            contentDiv.innerHTML = renderWhatWentWrong(matchData);
            break;
        case 'critical-moments':
            contentDiv.innerHTML = renderCriticalMoments(matchData);
            break;
        case 'items':
            contentDiv.innerHTML = renderItemAnalysis(matchData);
            break;
        case 'matchup':
            contentDiv.innerHTML = renderMatchupAnalysis(matchData);
            break;
        case 'deep-dive':
            contentDiv.innerHTML = renderDeepDive(matchData);
            break;
    }
}

function renderOverview(data) {
    let html = '<div class="stats-grid">';
    html += `<div class="stat-card"><div class="label">Match ID</div><div class="value">${data.match_id || 'N/A'}</div></div>`;
    
    if (data.structured_insights && data.structured_insights.key_statistics) {
        const stats = data.structured_insights.key_statistics;
        if (stats.combat) {
            stats.combat.forEach(stat => {
                html += `<div class="stat-card"><div class="label">${stat.label}</div><div class="value">${stat.value}</div></div>`;
            });
        }
    }
    
    html += '</div>';
    
    html += '<div class="deep-dive-content">';
    html += `<h2>Match Analysis</h2>`;
    html += `<p>${data.analysis || 'Analysis not available'}</p>`;
    
    if (data.suggestions && data.suggestions.length > 0) {
        html += `<h3 style="margin-top: 24px;">Suggestions</h3><ul>`;
        data.suggestions.forEach(s => {
            html += `<li>${s}</li>`;
        });
        html += `</ul>`;
    }
    
    html += `</div>`;
    return html;
}

function renderWhatWentWell(data) {
    if (!data.structured_insights || !data.structured_insights.what_went_well) {
        return '<p>Structured insights not available. Deep dive analysis may not have been generated.</p>';
    }

    let html = '<h2>✅ What Went Well</h2>';
    data.structured_insights.what_went_well.forEach(event => {
        html += `<div class="event-card positive">`;
        html += `<h3>${event.title}</h3>`;
        html += `<div class="description">${event.description}</div>`;
        if (event.impact) {
            html += `<div class="impact">${event.impact}</div>`;
        }
        if (event.data && event.data.length > 0) {
            html += `<div class="data">`;
            event.data.forEach(d => {
                html += `<span class="data-badge">${d}</span>`;
            });
            html += `</div>`;
        }
        html += `</div>`;
    });
    return html;
}

function renderWhatWentWrong(data) {
    if (!data.structured_insights || !data.structured_insights.what_went_wrong) {
        return '<p>Structured insights not available. Deep dive analysis may not have been generated.</p>';
    }

    let html = '<h2>❌ What Went Wrong</h2>';
    data.structured_insights.what_went_wrong.forEach(event => {
        html += `<div class="event-card negative">`;
        html += `<h3>${event.title}</h3>`;
        html += `<div class="description">${event.description}</div>`;
        if (event.impact) {
            html += `<div class="impact">${event.impact}</div>`;
        }
        if (event.data && event.data.length > 0) {
            html += `<div class="data">`;
            event.data.forEach(d => {
                html += `<span class="data-badge">${d}</span>`;
            });
            html += `</div>`;
        }
        html += `</div>`;
    });
    return html;
}

function renderCriticalMoments(data) {
    if (!data.structured_insights || !data.structured_insights.critical_moments) {
        return '<p>Structured insights not available. Deep dive analysis may not have been generated.</p>';
    }

    let html = '<h2>⚡ Critical Moments</h2>';
    data.structured_insights.critical_moments.forEach(moment => {
        html += `<div class="event-card critical">`;
        html += `<h3>${moment.title}</h3>`;
        html += `<div class="description">${moment.description}</div>`;
        if (moment.outcome) {
            html += `<div class="impact"><strong>Outcome:</strong> ${moment.outcome}</div>`;
        }
        if (moment.impact) {
            html += `<div class="impact">${moment.impact}</div>`;
        }
        if (moment.data && moment.data.length > 0) {
            html += `<div class="data">`;
            moment.data.forEach(d => {
                html += `<span class="data-badge">${d}</span>`;
            });
            html += `</div>`;
        }
        html += `</div>`;
    });
    return html;
}

function renderItemAnalysis(data) {
    if (!data.structured_insights || !data.structured_insights.item_analysis) {
        return '<p>Item analysis not available. Deep dive analysis may not have been generated.</p>';
    }

    const itemAnalysis = data.structured_insights.item_analysis;
    let html = '<h2>🛡️ Item Build Analysis</h2>';
    
    html += `<div class="deep-dive-content">`;
    if (itemAnalysis.timing_analysis) {
        html += `<h3>Timing Analysis</h3><p>${itemAnalysis.timing_analysis}</p>`;
    }
    if (itemAnalysis.opponent_matchup) {
        html += `<h3>Opponent Matchup</h3><p>${itemAnalysis.opponent_matchup}</p>`;
    }
    if (itemAnalysis.recommendations && itemAnalysis.recommendations.length > 0) {
        html += `<h3>Recommendations</h3><ul>`;
        itemAnalysis.recommendations.forEach(r => {
            html += `<li>${r}</li>`;
        });
        html += `</ul>`;
    }
    html += `</div>`;
    return html;
}

function renderMatchupAnalysis(data) {
    if (!data.structured_insights || !data.structured_insights.matchup_analysis) {
        return '<p>Matchup analysis not available. Deep dive analysis may not have been generated.</p>';
    }

    const matchup = data.structured_insights.matchup_analysis;
    let html = '<h2>⚔️ Matchup & Team Composition Analysis</h2>';
    
    html += `<div class="deep-dive-content">`;
    if (matchup.lane_matchup) {
        html += `<h3>Lane Matchup</h3><p>${matchup.lane_matchup}</p>`;
    }
    if (matchup.team_composition) {
        html += `<h3>Team Composition</h3><p>${matchup.team_composition}</p>`;
    }
    if (matchup.synergies && matchup.synergies.length > 0) {
        html += `<h3>Synergies</h3><ul>`;
        matchup.synergies.forEach(s => {
            html += `<li>${s}</li>`;
        });
        html += `</ul>`;
    }
    if (matchup.counters && matchup.counters.length > 0) {
        html += `<h3>Counters</h3><ul>`;
        matchup.counters.forEach(c => {
            html += `<li>${c}</li>`;
        });
        html += `</ul>`;
    }
    if (matchup.win_conditions && matchup.win_conditions.length > 0) {
        html += `<h3>Win Conditions</h3><ul>`;
        matchup.win_conditions.forEach(w => {
            html += `<li>${w}</li>`;
        });
        html += `</ul>`;
    }
    html += `</div>`;
    return html;
}

function renderDeepDive(data) {
    let html = '<h2>🔍 Deep Dive Analysis</h2>';
    html += `<div class="deep-dive-content">`;
    html += data.champion_deep_dive || 'Deep dive analysis not available. Please specify a champion_name or summoner_name.';
    html += `</div>`;
    return html;
}

// Allow Enter key to submit
document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('gameId').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            analyzeMatch();
        }
    });
});

