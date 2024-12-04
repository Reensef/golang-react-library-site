// import React from 'react'
import { Box, useDisclosure } from '@chakra-ui/react'
import SidebarContent from './components/SidebarContent'
import MobileNav from './components/MobileNav'
import TableContent from './components/TableContent'

import { useColorModeValue } from '@chakra-ui/react'

function App() {
  const { onOpen, onClose } = useDisclosure()

  return (
    <Box minH="100vh" bg={useColorModeValue('gray.100', 'gray.900')}>
      {/* Sidebar */}
      <SidebarContent onClose={onClose} display={{ base: 'none', md: 'block' }} />

      {/* Mobile Navigation */}
      <MobileNav onOpen={onOpen} />

      {/* Main Content Area */}
      <Box ml={{ base: 0, md: 60 }} p="4">
        <TableContent />
      </Box>
    </Box>
  )
}

export default App
