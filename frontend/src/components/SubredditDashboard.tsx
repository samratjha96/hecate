import React, { useState, useEffect, useRef } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

interface Subreddit {
  name: string;
  numberOfSubscribers: number;
}

const SubredditDashboard = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

  const [subreddits, setSubreddits] = useState<Subreddit[]>([]);
  const fetchedSubredditsRef = useRef(subreddits);
  const [newSubreddit, setNewSubreddit] = useState("");
  const [timeRange, setTimeRange] = useState("day");

  useEffect(() => {
    fetchedSubredditsRef.current = subreddits;
  });
  useEffect(() => {
    fetchSubreddits();
  }, [subreddits]);

  const fetchSubreddits = async () => {
    try {
      const response = await fetch("http://localhost:8000/subreddits/");
      const data = await response.json();
      if (fetchedSubredditsRef.current.length !== data.length) {
        setSubreddits(data);
      }
    } catch (error) {
      console.error("Error fetching subreddits:", error);
    }
  };

  const handleIngestion = async (subredditName: string, timeRange: string) => {
    console.log("Ingesting:", { subredditName, timeRange });
    const response = await fetch(`${apiUrl}/subreddits/ingest`, {
      method: "POST",
      body: JSON.stringify({
        subreddit: { name: subredditName, sortBy: timeRange },
      }),
    });
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    const data = await response.json();

    toast(`Ingesting data from r/${data.Name}`);
  };

  const handleAddSubreddit = async (e: React.FormEvent<EventTarget>) => {
    e.preventDefault();
    if (newSubreddit) {
      console.log(
        `Adding new subreddit: ${newSubreddit} with time range: ${timeRange}`,
      );
      setNewSubreddit("");
      const response = await fetch(`${apiUrl}/subreddits/ingest`, {
        method: "POST",
        body: JSON.stringify({
          subreddits: { name: newSubreddit, sortBy: timeRange },
        }),
      });
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const data = await response.json();
      const addedSubreddit: Subreddit = {
        name: data.Name,
        numberOfSubscribers: data.NumberOfSubscribers,
      };
      setSubreddits((oldArray) => [...oldArray, addedSubreddit]);
      toast(`Subscribed to r/${newSubreddit}`, {});
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Subscribed Subreddits</h1>

      {subreddits.length > 0 && (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Subscribers</TableHead>
              <TableHead>Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {subreddits.map((subreddit) => (
              <TableRow key={subreddit.name}>
                <TableCell>{subreddit.name}</TableCell>
                <TableCell>{subreddit.numberOfSubscribers}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      onClick={() => handleIngestion(subreddit.name, "day")}
                    >
                      Ingest Day
                    </Button>
                    <Button
                      variant="secondary"
                      onClick={() => handleIngestion(subreddit.name, "month")}
                    >
                      Ingest Month
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>Add New Subreddit</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAddSubreddit} className="space-y-4">
            <Input
              type="text"
              placeholder="Enter subreddit name"
              value={newSubreddit}
              onChange={(e) => setNewSubreddit(e.target.value)}
            />
            <Select value={timeRange} onValueChange={setTimeRange}>
              <SelectTrigger>
                <SelectValue placeholder="Select time range" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="day">Day</SelectItem>
                <SelectItem value="month">Month</SelectItem>
              </SelectContent>
            </Select>
          </form>
        </CardContent>
        <CardFooter>
          <Button type="submit" onClick={handleAddSubreddit}>
            Add Subreddit
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default SubredditDashboard;
