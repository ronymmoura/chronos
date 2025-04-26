import { Button } from "@kamalion/ui";
import { useQuery } from "@tanstack/react-query";
import { ProcessApi } from "../../api/process";
import { useEffect } from "react";

export function HomePage() {
  const { data: processList, isLoading, refetch } = useQuery({ queryKey: ["processList"], queryFn: ProcessApi.list });

  useEffect(() => {
    const refresher = setInterval(() => {
      refetch();
    }, 1000);

    return () => clearInterval(refresher);
  }, []);

  if (isLoading) {
    return <div className="p-5">Loading...</div>;
  }

  return (
    <div className="p-5">
      {processList?.length === 0 && <div className="flex h-full flex-col items-center justify-center">No items found. Please add a new process.</div>}

      <div>
        {processList?.map((process) => (
          <div
            key={process.id}
            className="mb-4 flex items-center justify-between rounded-lg bg-zinc-700 p-4 text-sm text-white shadow-md hover:bg-zinc-600"
          >
            <div>{process.running.toString()}</div>
            <div className="flex flex-col">
              <span className="font-bold">{process.name}</span>
              <span>{process.description}</span>
            </div>
            <div className="flex items-center gap-2">
              <Button.Root variant="accent" size="sm" onClick={() => alert(`Edit ${process.name}`)}>
                <Button.Content>Edit</Button.Content>
              </Button.Root>
              <Button.Root variant="danger" size="sm" onClick={() => alert(`Delete ${process.name}`)}>
                <Button.Content>Delete</Button.Content>
              </Button.Root>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
