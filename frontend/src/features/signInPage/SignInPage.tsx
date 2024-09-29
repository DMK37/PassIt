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
  Link,
} from "@chakra-ui/react";
import useUserStore from "../../stores/userStore";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { login } from "../../apiCalls/userServiceCalls";

const SignInPage: React.FC = () => {
  const setUser = useUserStore((state) => state.setUser);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSignIn = async () => {
    // Simulate authentication
    try {
      if (!email || !password) {
        setError("Email and password are required");
        return;
      }

      // Call the login function from the user service
      const response = await login({ email, password });

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
            Sign In
          </Heading>
          {error && <Text color="red.500">{error}</Text>}
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
            onClick={handleSignIn}
            py={6}
            fontSize="1.1rem"
            color="background"
            _hover={{ bgColor: "accent" }}
          >
            Sign In
          </Button>
          <Text mt={4} textAlign="center">
            Don't have an account?{" "}
            <Link
              as={RouterLink}
              to="/signup"
              color="primary"
              _hover={{ color: "accent" }}
            >
              Sign Up
            </Link>
          </Text>
        </VStack>
      </Box>
    </Container>
  );
};

export default SignInPage;
