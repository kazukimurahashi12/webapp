import PropTypes from "prop-types";
import NextLink from "next/link";
import { Box, Typography, Unstable_Grid2 as Grid, Button } from "@mui/material";
import { Logo } from "../../../src/components/logo";

export const Layout = (props) => {
  const { children } = props;

  return (
    <Box
      component="main"
      sx={{
        display: "flex",
        flex: "1 1 auto",
      }}
    >
      <Grid container>
        <Box
          component="header"
          sx={{
            left: 0,
            p: 3,
            position: "fixed",
            top: 0,
            width: "100%",
          }}
        >
          <Box
            component={NextLink}
            href="/auth/overview"
            sx={{
              display: "inline-flex",
              height: 32,
              width: 32,
            }}
          >
            <Logo />
          </Box>
        </Box>
        {children}
      </Grid>
    </Box>
  );
};

Layout.prototypes = {
  children: PropTypes.node,
};
