import {readFile, writeFile} from "fs";

const toString = (coord) => `${coord[0]}_${coord[1]}`;

readFile('adjacentList.json', function(err, data) { 
    const adjacentList = JSON.parse(data); 
    const visited = new Set();
    const m = new Map();

    for (const pair of adjacentList) {
        const a = toString(pair[0]);
        const b = toString(pair[1]);

        if (!m.has(a)) m.set(a, []);
        if (!m.has(b)) m.set(b, []);

        m.get(a).push(b);
        m.get(b).push(a);

        visited.add(a);
        visited.add(b); 
    }

    console.log(visited.size);

    const queue =  ['151.2111385_-33.8659373'];
    while (queue.length > 0) {
        // Get next node
        const currentNode = queue.shift();

        // Mark as visited
        visited.delete(currentNode);

        // Queue all the adjacent nodes of current node
        for (const nextNode of m.get(currentNode)) {
            if (visited.has(nextNode)) queue.push(nextNode)
        }
    }

    for (const example of visited) {
        for (const pair of adjacentList) {
            const a = toString(pair[0]);
            const b = toString(pair[1]);

            if (a === example || b === example) {
                const idx = adjacentList.findIndex(p => p == pair);
                adjacentList.splice(idx, 1);
            }
        }
    }

    writeFile('adjacentList2.json', JSON.stringify(adjacentList), () => {})
}); 
