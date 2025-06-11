    #include <iostream>
    #include <string>
    #include <vector>
    #include <chrono>
    #include <cstdlib>
    #include <cstring>
    #include <unistd.h> 
    #include <arpa/inet.h> 
    #include <sys/socket.h>
    #include <netinet/in.h> 

    #include "tuning.pb.h"

    uint64_t current_time_millis() {
        using namespace std::chrono;
        return duration_cast<milliseconds>(system_clock::now().time_since_epoch()).count();
    }

    int main(int argc, char* argv[]) {
        GOOGLE_PROTOBUF_VERIFY_VERSION;
        float fuzzValue = 0.0f;

            std::string input;
            if (!std::getline(std::cin, input)) {
                std::cerr << "Failed to read input from stdin\n";
                return 1;
            }
            fuzzValue = std::stof(input);

        protobuf_msgs::TuningState tuning;
        tuning.set_timestamp(current_time_millis());

        protobuf_msgs::TuningState::Parameter* param = tuning.add_dynamicparameters();
        protobuf_msgs::TuningState::Parameter::NumberParameter* numberParam = new protobuf_msgs::TuningState::Parameter::NumberParameter();
        numberParam->set_key("speed");
        numberParam->set_value(fuzzValue);
        param->set_allocated_number(numberParam);

        std::string data;
        if (!tuning.SerializeToString(&data)) {
            std::cerr << "Failed to serialize protobuf message\n";
            return 1;
        }

        int sock = socket(AF_INET, SOCK_STREAM, 0);
        if (sock < 0) {
            perror("socket");
            return 1;
        }

        sockaddr_in serv_addr{};
        serv_addr.sin_family = AF_INET;
        serv_addr.sin_port = htons(9000);
        if (inet_pton(AF_INET, "192.168.0.146", &serv_addr.sin_addr) <= 0) {
            std::cerr << "Invalid address/ Address not supported\n";
            close(sock);
            return 1;
        }

        if (connect(sock, (sockaddr*)&serv_addr, sizeof(serv_addr)) < 0) {
            perror("connect");
            close(sock);
            return 1;
        }

        ssize_t total_sent = 0;
        ssize_t to_send = data.size();
        const char* buffer = data.data();
        while (total_sent < to_send) {
            ssize_t sent = send(sock, buffer + total_sent, to_send - total_sent, 0);
            if (sent <= 0) {
                perror("send");
                close(sock);
                return 1;
            }
            total_sent += sent;
        }
        
        close(sock);
        google::protobuf::ShutdownProtobufLibrary();
        std::cout << "TuningState sent to transceiver successfully.\n";
        return 0;
    }
