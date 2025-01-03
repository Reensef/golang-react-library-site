import {
  Flex,
  IconButton,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  FlexProps,
  useColorModeValue,
} from "@chakra-ui/react";
import { FiMenu } from "react-icons/fi";
import { FaUserCircle } from "react-icons/fa";

interface MobileProps extends FlexProps {
  onOpen: () => void;
}

const MobileNav = ({ onOpen, ...rest }: MobileProps) => {
  const handleSignOut = () => {
    localStorage.removeItem("token");
    window.location.href = "/loginpage";
    window.history.replaceState(null, "", "/loginpage");
  };

  return (
    <Flex
      ml={{ base: 0, md: 60 }}
      px={{ base: 4, md: 4 }}
      height="20"
      alignItems="center"
      bg={useColorModeValue("white", "gray.900")}
      borderBottomColor={useColorModeValue("gray.200", "gray.700")}
      justifyContent={{ base: "space-between", md: "flex-end" }}
      {...rest}
    >
      <IconButton
        display={{ base: "flex", md: "none" }}
        onClick={onOpen}
        variant="outline"
        aria-label="open menu"
        icon={<FiMenu />}
      />
      <Flex alignItems="center">
        <Menu>
          <MenuButton
            as={IconButton}
            aria-label="Actions"
            size={"lg"}
            icon={<FaUserCircle />}
            variant="ghost"
          />
          <MenuList onClick={() => handleSignOut()}>
            <MenuItem>Sign out</MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Flex>
  );
};

export default MobileNav;
