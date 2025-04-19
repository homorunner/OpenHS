// Game state and UI variables
let gameState = null;
let selectedCard = null;
let selectedMinion = null;
let isAttacking = false;

// Card type icons
const cardIcons = {
    Minion: "☺",
    Spell: "✧",
    Weapon: "⚔",
    Hero: "♛",
    "Hero Power": "⚡"
};

// Initialize the game
document.addEventListener('DOMContentLoaded', () => {
    // Fetch initial game state
    fetchGameState();
    
    // Set up event listeners
    document.getElementById('end-turn').addEventListener('click', endTurn);
    
    // Poll for game state updates every 3 seconds (for demo purposes)
    // In a real game, you'd use websockets for real-time updates
    // setInterval(fetchGameState, 3000);
});

// Fetch game state from the server
function fetchGameState() {
    fetch('/api/game')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            gameState = data;
            updateUI();
        })
        .catch(error => {
            console.error('Error fetching game state:', error);
            logMessage('Error: Could not connect to the game server.');
        });
}

// Update the UI based on the game state
function updateUI() {
    if (!gameState) return;
    
    // Update turn and phase info
    document.getElementById('turn-info').textContent = `Turn: ${gameState.currentTurn}`;
    document.getElementById('phase-info').textContent = `Phase: ${gameState.phase}`;
    
    // Get player and opponent based on current player index
    const currentPlayerIdx = gameState.currentPlayerIndex;
    const player = gameState.players[currentPlayerIdx];
    const opponent = gameState.players[1 - currentPlayerIdx];
    
    // Update player hero
    document.getElementById('player-hero-name').textContent = player.hero.name;
    document.getElementById('player-hero-health').textContent = `${player.hero.health} HP`;
    
    // Update opponent hero
    document.getElementById('opponent-hero-name').textContent = opponent.hero.name;
    document.getElementById('opponent-hero-health').textContent = `${opponent.hero.health} HP`;
    
    // Update mana displays
    document.getElementById('player-mana').textContent = `${player.mana}/${player.totalMana}`;
    document.getElementById('opponent-mana').textContent = `${opponent.mana}/${opponent.totalMana}`;
    
    // Update player's hand
    updateHand('player-hand', player.hand);
    
    // For opponent's hand, we show card backs only
    updateOpponentHand('opponent-hand', opponent.hand.length);
    
    // Update player's field
    updateField('player-field', player.field, true);
    
    // Update opponent's field
    updateField('opponent-field', opponent.field, false);
    
    // Reset selection states
    selectedCard = null;
    selectedMinion = null;
    isAttacking = false;
}

// Update the player's hand display
function updateHand(containerId, cards) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';
    
    cards.forEach((card, index) => {
        const cardElement = createCardElement(card, index);
        container.appendChild(cardElement);
    });
}

// Update the opponent's hand display (card backs only)
function updateOpponentHand(containerId, cardCount) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';
    
    for (let i = 0; i < cardCount; i++) {
        const cardElement = document.importNode(document.getElementById('card-template').content, true).querySelector('.card');
        
        // Set card back appearance
        const cardIcon = cardElement.querySelector('.card-icon');
        cardIcon.textContent = "?";
        
        container.appendChild(cardElement);
    }
}

// Update the field display
function updateField(containerId, minions, isPlayerField) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';
    
    minions.forEach((minion, index) => {
        const minionElement = createMinionElement(minion, index, isPlayerField);
        container.appendChild(minionElement);
    });
}

// Create a card element
function createCardElement(card, index) {
    const cardElement = document.importNode(document.getElementById('card-template').content, true).querySelector('.card');
    
    cardElement.dataset.index = index;
    
    const cardName = cardElement.querySelector('.card-name');
    cardName.textContent = card.name;
    
    const cardCost = cardElement.querySelector('.card-cost');
    cardCost.textContent = card.cost;
    
    const cardDescription = cardElement.querySelector('.card-description');
    cardDescription.textContent = card.description || '';
    
    const cardIcon = cardElement.querySelector('.card-icon');
    cardIcon.textContent = cardIcons[card.type] || cardIcons.Minion;
    
    // Only show attack/health for minions and weapons
    if (card.type === 'Minion' || card.type === 'Weapon') {
        const cardAttack = cardElement.querySelector('.card-attack');
        cardAttack.textContent = card.attack;
        
        const cardHealth = cardElement.querySelector('.card-health');
        cardHealth.textContent = card.health;
    }
    
    // Add click event for player's cards
    cardElement.addEventListener('click', () => handleCardClick(index));
    
    return cardElement;
}

// Create a minion element
function createMinionElement(minion, index, isPlayerMinion) {
    const minionElement = document.importNode(document.getElementById('minion-template').content, true).querySelector('.minion');
    
    minionElement.dataset.index = index;
    
    const minionName = minionElement.querySelector('.minion-name');
    minionName.textContent = minion.name;
    
    const minionAttack = minionElement.querySelector('.minion-attack');
    minionAttack.textContent = minion.attack;
    
    const minionHealth = minionElement.querySelector('.minion-health');
    minionHealth.textContent = minion.health;
    
    const minionTags = minionElement.querySelector('.minion-tags');
    if (minion.tags && minion.tags.length > 0) {
        minionTags.textContent = minion.tags.join(', ');
    }
    
    // Add can-attack class if the minion can attack
    if (isPlayerMinion && minion.canAttack) {
        minionElement.classList.add('can-attack');
    }
    
    // Add click event
    if (isPlayerMinion) {
        minionElement.addEventListener('click', () => handlePlayerMinionClick(index));
    } else {
        minionElement.addEventListener('click', () => handleOpponentMinionClick(index));
    }
    
    return minionElement;
}

// Handle click on a card in hand
function handleCardClick(index) {
    if (isAttacking) return; // Don't allow playing cards during attack selection
    
    const cardElement = document.querySelector(`#player-hand .card[data-index="${index}"]`);
    
    // Deselect if already selected
    if (selectedCard === index) {
        cardElement.style.transform = '';
        selectedCard = null;
        return;
    }
    
    // Deselect previous card if any
    if (selectedCard !== null) {
        const prevCardElement = document.querySelector(`#player-hand .card[data-index="${selectedCard}"]`);
        if (prevCardElement) {
            prevCardElement.style.transform = '';
        }
    }
    
    // Select new card
    selectedCard = index;
    cardElement.style.transform = 'translateY(-20px)';
    
    // For minions, we need to choose a position
    const card = gameState.players[gameState.currentPlayerIndex].hand[index];
    if (card.type === 'Minion') {
        logMessage('Select a position on the board to play this minion.');
    } else {
        // For non-minions, play immediately
        playCard(index);
    }
}

// Handle click on player's minion
function handlePlayerMinionClick(index) {
    const minion = gameState.players[gameState.currentPlayerIndex].field[index];
    const minionElement = document.querySelector(`#player-field .minion[data-index="${index}"]`);
    
    // If we're placing a minion from hand
    if (selectedCard !== null) {
        // Calculate position to play the card
        const fieldSize = gameState.players[gameState.currentPlayerIndex].field.length;
        const position = Math.min(index + 1, fieldSize);
        playCard(selectedCard, position);
        return;
    }
    
    // Initiate attack if the minion can attack
    if (minion.canAttack) {
        // Deselect if already selected
        if (selectedMinion === index && isAttacking) {
            minionElement.classList.remove('selected');
            selectedMinion = null;
            isAttacking = false;
            logMessage('Attack cancelled.');
            return;
        }
        
        // Deselect previous minion if any
        if (selectedMinion !== null) {
            const prevMinionElement = document.querySelector(`#player-field .minion[data-index="${selectedMinion}"]`);
            if (prevMinionElement) {
                prevMinionElement.classList.remove('selected');
            }
        }
        
        // Select new minion for attacking
        selectedMinion = index;
        isAttacking = true;
        minionElement.classList.add('selected');
        logMessage('Select a target to attack.');
    }
}

// Handle click on opponent's minion
function handleOpponentMinionClick(index) {
    // If we're not attacking, do nothing
    if (!isAttacking || selectedMinion === null) {
        return;
    }
    
    // Attack the selected opponent minion
    attack(selectedMinion, index);
}

// Play a card from hand
function playCard(cardIndex, position = null) {
    // Reset selections
    selectedCard = null;
    
    // Prepare action data
    const actionData = {
        type: 'playCard',
        cardIndex: cardIndex,
        position: position !== null ? position : -1
    };
    
    // Send action to server
    sendAction(actionData);
}

// Attack with a minion
function attack(attackerIndex, targetIndex) {
    // Reset selections
    selectedMinion = null;
    isAttacking = false;
    
    // Remove selected class
    const minionElement = document.querySelector(`#player-field .minion.selected`);
    if (minionElement) {
        minionElement.classList.remove('selected');
    }
    
    // Prepare action data
    const actionData = {
        type: 'attack',
        cardIndex: attackerIndex,
        target: targetIndex
    };
    
    // Send action to server
    sendAction(actionData);
}

// End the current player's turn
function endTurn() {
    // Reset selections
    selectedCard = null;
    selectedMinion = null;
    isAttacking = false;
    
    // Prepare action data
    const actionData = {
        type: 'endTurn'
    };
    
    // Send action to server
    sendAction(actionData);
}

// Send an action to the server
function sendAction(actionData) {
    fetch('/api/action', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(actionData)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                throw new Error(text || 'Action failed');
            });
        }
        return response.json();
    })
    .then(data => {
        gameState = data;
        updateUI();
        logMessage(`Action performed successfully.`);
    })
    .catch(error => {
        console.error('Error:', error);
        logMessage(`Error: ${error.message}`);
    });
}

// Log a message to the game log
function logMessage(message) {
    const gameLog = document.getElementById('game-log');
    const logEntry = document.createElement('p');
    logEntry.textContent = message;
    gameLog.appendChild(logEntry);
    
    // Auto-scroll to bottom
    gameLog.scrollTop = gameLog.scrollHeight;
} 