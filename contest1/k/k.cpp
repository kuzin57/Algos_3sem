#include <iostream>
#include <map>
#include <memory>
#include <vector>
#include <algorithm>
#include <string.h>

struct node {
    int next[26];
    int link;
    int len;
    node() {
        for(int i = 0; i < 26; ++i) {
            next[i] = -1;
        }
    }

    node(int len, int link): len(len), link(link) {
        for(int i = 0; i < 26; ++i) {
            next[i] = -1;
        }
    }

    ~node() = default;
};

const int MAXSIZE = 100000;
node nodes[MAXSIZE * 2];

struct Automata {
    int sz;
    int last;

    Automata();
    void Add(char c);
    void AddWord(const std::string& line);
    bool Find(const std::string& word);

    ~Automata() {
        
    }
};

Automata::Automata() {
    sz = 0;
    last = 0;
    nodes[sz++] = node(0, -1);
}

void Automata::Add(char c) {
    bool exit = false;
    nodes[sz++].len = nodes[last].len + 1;
    nodes[sz-1].link = -1;
    auto cur_node = last;
    last = sz-1;
    while (cur_node != -1) {
        auto to = nodes[cur_node].next[c-'a'];
        if (to == -1) {
            nodes[cur_node].next[c-'a'] = sz-1;
        } else {
            if (nodes[to].len == nodes[cur_node].len + 1) {
                nodes[last].link = to;
            } else {
                nodes[sz++].len = nodes[cur_node].len + 1;
                nodes[sz-1].link = -1;
                for (int i = 0; i < 26; ++i) {
                    nodes[sz-1].next[i] = nodes[to].next[i];
                }
    
                nodes[sz-1].link = nodes[to].link;
                nodes[to].link = sz-1;
                nodes[last].link = sz-1;
                auto suff_link = cur_node;
                bool flag = false;
                while (suff_link != -1) {
                    auto n = nodes[suff_link].next[c-'a'];
                    if (n == to) {
                        nodes[suff_link].next[c-'a'] = sz-1;
                    } else {
                        break;
                    }
                    suff_link = nodes[suff_link].link;
                }
            }
            break;
        }
        cur_node = nodes[cur_node].link;
    }
    if (cur_node == -1) {
        nodes[last].link = 0;
    }
}


void Automata::AddWord(const std::string& line) {
    for (int i = 0; i < line.size(); ++i) {
        if (line[i] >= 'A' && line[i] <= 'Z') {
            Add('a' + line[i] - 'A');
        } else {
            Add(line[i]);
        }
    }
}

bool Automata::Find(const std::string& word) {
    int cur_node = 0;
    int cur_index = 0;
    while (cur_index < word.size()) {
        int to;
        if (word[cur_index] >= 'A' && word[cur_index] <= 'Z') {
            to = nodes[cur_node].next[word[cur_index] - 'A'];
        } else {
            to = nodes[cur_node].next[word[cur_index] - 'a'];
        }
        if (to == -1) {
            return false;
        }
        ++cur_index;
        cur_node = to;
    }
    return true;
}

int main() {
    std::ios_base::sync_with_stdio(false);
    std::cin.tie(NULL);
    Automata automata;
    std::string operation = "";
    std::string line = "";
    while(std::cin >> operation) {
        std::cin >> line;
        if (operation == "A") {
            automata.AddWord(line);
        } else {
            bool res = automata.Find(line);
            switch (res) {
            case true:
                std::cout << "YES" << '\n';
                break;
            case false:
                std::cout << "NO" << '\n';
                break;
            }
        }
    }
}