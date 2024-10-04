import React from "react";
import { useNavigate } from "react-router-dom";
import useUserStore, { isTokenExpired } from "../stores/userStore";
import { Box } from "@chakra-ui/react";

interface ProtectedRouteProps {
  element: React.ReactElement;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ element }) => {
    const { token, logout } = useUserStore((state) => ({
      token: state.token,
      logout: state.logout,
    }));
    const navigate = useNavigate();
  
    if (isTokenExpired(token)) {
      logout();
      navigate("/signin");
      return <Box></Box>;
    }
  
    return element;
  };

export default ProtectedRoute;