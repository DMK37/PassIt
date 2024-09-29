import React, { useState } from "react";
import {
  Box,
  Button,
  Container,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Heading,
  Text,
} from "@chakra-ui/react";
import useUserStore from "../../stores/userStore";
import { useNavigate } from "react-router-dom";
import { signUp } from "../../apiCalls/userServiceCalls";

const SignUpPage: React.FC = () => {
  const setUser = useUserStore((state) => state.setUser);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [username, setUsername] = useState("");
  const navigate = useNavigate();

  const handleSignUp = async () => {
    try {
      if (!email || !password || !firstName || !lastName || !username) {
        setError("All fields are required");
        return;
      }

      const response = await signUp({
        username,
        email,
        password,
        first_name: firstName,
        last_name: lastName,
      });
      setUser(response.user, response.token);
      navigate("/");
    } catch (error) {
      setError("Invalid email or password");
    }
  };

  return (
    <Container maxW="lg" mt={20}>
      <Box p={6} bgColor="transparent" borderRadius="15" color="text">
        <VStack spacing={4} align="stretch">
          <Heading size="lg" textAlign="center">
            Sign Up
          </Heading>
          {error && <Text color="red.500">{error}</Text>}
          <FormControl id="username">
            <FormLabel>Username</FormLabel>
            <Input
              placeholder="Username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              _placeholder={{ color: "text" }}
              borderWidth={1}
              borderColor="grey"
              fontSize="1.1rem"
              py="7"
              _hover={{ borderColor: "primary" }}
              sx={{
                ":-webkit-autofill": {
                  WebkitBoxShadow: "0 0 0 30px white inset !important",
                  WebkitTextFillColor: "black !important",
                },
              }}
            />
          </FormControl>
          <FormControl id="firstName">
            <FormLabel>First Name</FormLabel>
            <Input
              placeholder="First Name"
              type="text"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              _placeholder={{ color: "text" }}
              borderWidth={1}
              borderColor="grey"
              fontSize="1.1rem"
              py="7"
              _hover={{ borderColor: "primary" }}
              sx={{
                ":-webkit-autofill": {
                  WebkitBoxShadow: "0 0 0 30px white inset !important",
                  WebkitTextFillColor: "black !important",
                },
              }}
            />
          </FormControl>
          <FormControl id="lastName">
            <FormLabel>Last Name</FormLabel>
            <Input
              placeholder="Last Name"
              type="text"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              _placeholder={{ color: "text" }}
              borderWidth={1}
              borderColor="grey"
              fontSize="1.1rem"
              py="7"
              _hover={{ borderColor: "primary" }}
              sx={{
                ":-webkit-autofill": {
                  WebkitBoxShadow: "0 0 0 30px white inset !important",
                  WebkitTextFillColor: "black !important",
                },
              }}
            />
          </FormControl>

          <FormControl id="email">
            <FormLabel>Email address</FormLabel>
            <Input
              placeholder="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              _placeholder={{ color: "text" }}
              borderWidth={1}
              borderColor="grey"
              fontSize="1.1rem"
              py="7"
              _hover={{ borderColor: "primary" }}
              sx={{
                ":-webkit-autofill": {
                  WebkitBoxShadow: "0 0 0 30px white inset !important",
                  WebkitTextFillColor: "black !important",
                },
              }}
            />
          </FormControl>
          <FormControl id="password">
            <FormLabel fontSize="1.1rem">Password</FormLabel>
            <Input
              placeholder="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              _placeholder={{ color: "text" }}
              borderWidth={1}
              borderColor="grey"
              fontSize="1.1rem"
              py="7"
              _hover={{ borderColor: "primary" }}
              sx={{
                ":-webkit-autofill": {
                  WebkitBoxShadow: "0 0 0 30px white inset !important",
                  WebkitTextFillColor: "black !important",
                },
              }}
            />
          </FormControl>
          <Button
            bgColor="primary"
            onClick={handleSignUp}
            py={6}
            fontSize="1.1rem"
            color="background"
            _hover={{ bgColor: "accent" }}
          >
            Sign Up
          </Button>
        </VStack>
      </Box>
    </Container>
  );
};

export default SignUpPage;
