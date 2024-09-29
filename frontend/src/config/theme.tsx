import { extendTheme } from '@chakra-ui/react';

// Define custom colors
const colors = {
    primary: "#2f27ce",
    secondary: "#dedcff",
    background: "#fbfbfe",
    text: "#050315",
    accent: "#433bff"
}

// Define custom fonts
const fonts = {
  heading: 'Roboto, sans-serif',
  body: 'Roboto, sans-serif',
};

// Define custom styles
const styles = {
  global: {
    body: {
        
      bg: 'background',
    },
  },
};

// Extend the default theme
const theme = extendTheme({ colors, fonts, styles });

export default theme;