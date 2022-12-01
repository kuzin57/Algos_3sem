#include <vector>
#include <memory>
#include <iostream>
#include <algorithm>

struct Vector {
    int64_t x, y;

    Vector();
    Vector(int64_t, int64_t);

    friend std::istream& operator>>(std::istream&, Vector&);
    friend std::ostream& operator<<(std::ostream&, Vector&);
};

Vector::Vector() = default;

Vector::Vector(int64_t x, int64_t y): x(x), y(y) {}

int64_t operator*(const Vector& first_vector, const Vector& second_vector) {
    return first_vector.x*second_vector.y - first_vector.y*second_vector.x;
}

Vector operator-(const Vector& first_vector, const Vector& second_vector) {
    return Vector(first_vector.x - second_vector.x, first_vector.y - second_vector.y);
}

bool operator<(const Vector& first_vector, const Vector& second_vector) {
    return first_vector.x < second_vector.x ||
        (first_vector.x == second_vector.x && first_vector.y < second_vector.y);
}

std::istream& operator>>(std::istream& in, Vector& vect) {
    in >> vect.x >> vect.y;
    return in;
}

std::ostream& operator<<(std::ostream& out, Vector& vect) {
    out << vect.x << " " << vect.y << '\n';
    return out;
}


struct Polygon {
    std::vector<Vector> vertices;
    Polygon(int);

    friend std::istream& operator>>(std::istream& in, Polygon& poly);
    friend std::ostream& operator<<(std::ostream& out, Polygon& poly);

    void convexHull(std::vector<Vector>&, std::vector<Vector>&);
};

Polygon::Polygon(int N) {
    vertices.resize(N);
}

std::istream& operator>>(std::istream& in, Polygon& poly) {
    for (size_t i = 0; i < poly.vertices.size(); ++i) {
        in >> poly.vertices[i];
    }
    return in;
}

std::ostream& operator<<(std::ostream& out, Polygon& poly) {
    for (size_t i = 0; i < poly.vertices.size(); ++i) {
        out << poly.vertices[i];
    }
    return out;
}

void Polygon::convexHull(std::vector<Vector>& upper_part, std::vector<Vector>& lower_part) {
    std::sort(vertices.begin(), vertices.end());

    upper_part.push_back(vertices[0]);
    upper_part.push_back(vertices[1]);
    lower_part.push_back(vertices[0]);
    lower_part.push_back(vertices[1]);

    auto process_vertex = [&](std::vector<Vector>& part, int is_upper, int index) {
        while (part.size() > 1 &&
            (part[part.size() - 2] - part[part.size() - 1]) * (vertices[index] - part[part.size() - 1]) * is_upper <= 0) {
            part.pop_back();
        }
    };

    for (size_t i = 2; i < vertices.size(); ++i) {
        process_vertex(upper_part, 1, i);
        upper_part.push_back(vertices[i]);
        process_vertex(lower_part, -1, i);
        lower_part.push_back(vertices[i]);
    }
}

int main() {
    std::ios_base::sync_with_stdio(false);
	std::cin.sync_with_stdio(false);
	std::cout.sync_with_stdio(false);

    int N;
    std::cin >> N;

    Polygon polygon(N);
    std::cin >> polygon;

    std::vector<Vector> upper_part, lower_part;

    polygon.convexHull(upper_part, lower_part);
    std::cout << (upper_part.size() + lower_part.size() - 2) << '\n';
    for (size_t i = 0; i < upper_part.size(); ++i) {
        std::cout << upper_part[i];
    }
    for (int i = lower_part.size() - 2; i >= 1; --i) {
        std::cout << lower_part[i];
    }
}
