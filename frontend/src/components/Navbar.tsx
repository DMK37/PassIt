import React from "react";
import {
  Box,
  Flex,
  IconButton,
  Input,
  Button,
  Spacer,
  InputGroup,
  InputRightElement,
} from "@chakra-ui/react";
import { VscAccount } from "react-icons/vsc";
import { SearchIcon } from "@chakra-ui/icons";
import { IoIosFootball } from "react-icons/io";
import { useNavigate } from "react-router-dom";
import useUserStore from "../stores/userStore";

const Navbar: React.FC = () => {
  const user = useUserStore((state) => state.user);
  const signout = useUserStore((state) => state.logout);
  const navigate = useNavigate();
  return (
    <Box
      px={4}
      py={2}
      bg="background"
      color="white"
      position="fixed"
      top={0}
      left={0}
      right={0}
      zIndex={1000}
      borderBottom="1px"
      borderBottomColor="gray.200"
    >
      <Flex alignItems="center">
        {/* Icon on the left side */}
        <Spacer />

        <IconButton
          aria-label="Home"
          icon={<IoIosFootball />}
          variant="ghost"
          color="primary"
          fontSize="32px"
          mr={4}
          _hover={{ color: "accent" }}
          _active={{ bg: "transparent" }}
          onClick={() => navigate("/")}
        />
        <Button
          fontSize="1.1rem"
          variant="ghost"
          color="text"
          fontWeight="normal"
          _hover={{ color: "gray.600" }}
          _active={{ bg: "transparent" }}
        >
          Apply
        </Button>
        {/* Spacer to push the search bar to the center */}
        <Spacer />

        {/* Search bar in the center */}
        <InputGroup maxW="400px" mx="auto">
          <Input
            fontSize="1.1rem"
            placeholder="Search..."
            variant="filled"
            bg="secondary"
            _placeholder={{ color: "text" }}
            color="text"
            _hover={{ bg: "secondary" }}
            _focus={{ bg: "secondary" }}
            borderRadius={20}
          />
          <InputRightElement>
            <SearchIcon color="text" />
          </InputRightElement>
        </InputGroup>

        {/* Spacer to push the buttons to the right */}
        <Spacer />

        {/* Profile and Sign Out buttons on the right side */}
        {user && (
          <IconButton
            variant="ghost"
            color="text"
            mr={2}
            aria-label="Account"
            icon={<VscAccount />}
            fontSize="30px"
            _hover={{ color: "gray.600" }}
            _active={{ bg: "transparent" }}
            onClick={() => navigate("/profile")}
          />
        )}
        {user && (
          <Button
            variant="ghost"
            color="text"
            fontWeight="normal"
            fontSize="1.1rem"
            _hover={{ color: "gray.600" }}
            _active={{ bg: "transparent" }}
            onClick={() => {
              signout();
              navigate("/");
            }}
          >
            Sign Out
          </Button>
        )}
        {!user && (
          <Button
            variant="ghost"
            color="text"
            fontWeight="normal"
            fontSize="1.1rem"
            _hover={{ color: "gray.600" }}
            _active={{ bg: "transparent" }}
            onClick={() => navigate("/signin")}
          >
            Sign In
          </Button>
        )}

        <Spacer />
      </Flex>
    </Box>
  );
};

export default Navbar;
