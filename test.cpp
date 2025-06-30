#include <iostream>
#include <fstream>
#include <vector>
#include <cstring>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>

int main() {
    const char* filepath = "crashes/crash_1751284192.bin";
    std::ifstream file(filepath, std::ios::binary);
    if (!file) {
        std::cerr << "Failed to open file: " << filepath << "\n";
        return 1;
    }

    std::vector<char> buffer((std::istreambuf_iterator<char>(file)),
                              std::istreambuf_iterator<char>());

    file.close();
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) {
        perror("socket");
        return 1;
    }

    sockaddr_in serv_addr{};
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(9000);
    if (inet_pton(AF_INET, "192.168.0.146", &serv_addr.sin_addr) <= 0) {
        std::cerr << "Invalid address\n";
        close(sock);
        return 1;
    }

    if (connect(sock, (sockaddr*)&serv_addr, sizeof(serv_addr)) < 0) {
        perror("connect");
        close(sock);
        return 1;
    }

    ssize_t total_sent = 0;
    ssize_t to_send = buffer.size();
    const char* data = buffer.data();

    while (total_sent < to_send) {
        ssize_t sent = send(sock, data + total_sent, to_send - total_sent, 0);
        if (sent <= 0) {
            perror("send");
            close(sock);
            return 1;
        }
        total_sent += sent;
    }

    close(sock);
    return 0;
}
