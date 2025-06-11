    #include <iostream>
    #include <string>
    #include <vector>
    #include <chrono>
    #include <cstdlib>
    #include <cstring>
    #include <unistd.h>     // read, close
    #include <arpa/inet.h>  // inet_pton
    #include <sys/socket.h> // socket, connect
    #include <netinet/in.h> // sockaddr_in

    #include "tuning.pb.h" // Generated protobuf header from tuning.proto

    // Helper: get current time in milliseconds since epoch
    uint64_t current_time_millis() {
        using namespace std::chrono;
        return duration_cast<milliseconds>(system_clock::now().time_since_epoch()).count();
    }

    int main(int argc, char* argv[]) {
        GOOGLE_PROTOBUF_VERIFY_VERSION;

        // 1. Read float value from argv or stdin
        float fuzzValue = 0.0f;
        if (argc > 1) {
            fuzzValue = std::stof(argv[1]);
        } else {
            // Read line from stdin
            std::string input;
            if (!std::getline(std::cin, input)) {
                std::cerr << "Failed to read input from stdin\n";
                return 1;
            }
            fuzzValue = std::stof(input);
        }

        // 2. Create protobuf message
        protobuf_msgs::TuningState tuning;
        tuning.set_timestamp(current_time_millis());

        // Add one dynamic parameter
        protobuf_msgs::TuningState::Parameter* param = tuning.add_dynamicparameters();
        protobuf_msgs::TuningState::Parameter::NumberParameter* numberParam = new protobuf_msgs::TuningState::Parameter::NumberParameter();
        numberParam->set_key("speed");
        numberParam->set_value(fuzzValue);
        param->set_allocated_number(numberParam);

        // 3. Serialize to string
        std::string data;
        if (!tuning.SerializeToString(&data)) {
            std::cerr << "Failed to serialize protobuf message\n";
            return 1;
        }

        // 4. Setup TCP connection to 192.168.0.146:9000
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

        // 5. Send serialized protobuf bytes directly (no size prefix in your Go code)
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
