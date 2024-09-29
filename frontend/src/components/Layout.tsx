
import { Box } from "@chakra-ui/react";
import Navbar from "./Navbar";
import { Outlet } from "react-router-dom";

const Layout: React.FC = () => {
    return (
        <Box minHeight="100vh" display="flex" flexDirection="column">
            <Navbar />
            <Box
                flex="1"
                display="flex"
                flexDirection="column"
                alignItems="center"
                width="100%"
                mt={8} 
            >
                <Outlet />
            </Box>
        </Box>
    );
};

export default Layout;