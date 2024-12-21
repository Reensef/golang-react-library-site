import { useEffect, useState } from "react";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
} from "@chakra-ui/react";
import axios from "axios";

interface ActionData {
  id: number;
  userId: number;
  userName: string;
  fileId: number;
  fileName: string;
  actionName: string;
  actionDateTime: string;
}

type ResponsActionData = {
  data: {
    id: number;
    user_id: number;
    user_name: string;
    files_id: number;
    file_name: string;
    action_name: string;
    created_at?: string;
  }[];
};

export const useFetchFilesActions = () => {
  const [actions, setActions] = useState<ActionData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchData = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await axios.get<ResponsActionData>(
        "/api/v1/files/actions_log",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response || !response.data || !Array.isArray(response.data.data)) {
        throw new Error("Invalid response format");
      }

      const data = response.data.data;
      const formattedData = data.map((item) => ({
        id: item.id,
        userId: item.user_id,
        userName: item.user_name,
        fileId: item.files_id,
        fileName: item.file_name,
        actionName: item.action_name,
        actionDateTime: item.created_at
          ? new Date(item.created_at).toLocaleString()
          : "",
      }));

      setActions(formattedData);
      setIsLoading(false);
    } catch (error) {
      setError(error as Error);
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  return { actions, isLoading, error };
};

const FilesActionsTable = () => {
  const { actions, isLoading, error } = useFetchFilesActions();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th width="30%">User</Th>
            <Th width="50%">Action</Th>
            <Th width="20%">Date & Time</Th>
          </Tr>
        </Thead>
        <Tbody>
          {actions.map((data) => (
            <Tr key={data.id}>
              <Td>{data.userName}</Td>
              <Td>{`${data.actionName}: ${data.fileName}`}</Td>
              <Td>{data.actionDateTime}</Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FilesActionsTable;
