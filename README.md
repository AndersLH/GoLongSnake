# GoLongSnake - IMT4306 DS project

## Description
This is a game made with Golang and Vue.js. It is a mix of both the game "Snake" and the game "Splatoon", as the user has to move around in a snake-like fashion, with the objective of coloring the entire board. Once the board has been filled up with colors and there are no more empty squares, the game ends, the score is tallied and the highest score wins. Your own snake color is always **cyan**. The board does not have wall collisions, instead you teleport through, and each player has a collision on their head, and only their head. Use the arrow keys to move around. 

## How to run
### Installation requirements
To run this locally on your machine, you will need Golang and NPM. 
- Golang: https://go.dev/doc/install
- NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm

### Build and run backend
- Open a command-line interface and navigate into the **backend** folder and run **go run .\cmd\main.go**

### Build and run frontend
- Open a second command-line interface and navigate into the **UI** folder with **cd ui**
- Run **npm i** to install required dependencies
- Run **npm run serve** for a development
