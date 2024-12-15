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
import { API_BASE_URL } from "@/lib/utils";
import SubredditPosts from "./SubredditPosts";

interface Subreddit {
  name: string;
  numberOfSubscribers: number;
}

const SubredditDashboard = () => {
  const [subreddits, setSubreddits] = useState<Subreddit[]>([]);
  const fetchedSubredditsRef = useRef(subreddits);
  const [newSubreddit, setNewSubreddit] = useState("");
  const [timeRange, setTimeRange] = useState("day");
  const [selectedSubreddit, setSelectedSubreddit] = useState<string>("");

  useEffect(() => {
    fetchedSubredditsRef.current = subreddits;
    // Set the first subreddit as selected when the list is loaded
    if (subreddits.length > 0 && !selectedSubreddit) {
      setSelectedSubreddit(subreddits[0].name);
    }
  }, [subreddits]);

  useEffect(() => {
    fetchSubreddits();
  }, []);

  const fetchSubreddits = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/subreddits/`);
      const data = await response.json();
      if (fetchedSubredditsRef.current.length !== data.length) {
        setSubreddits(data);
      }
    } catch (error) {
      console.error("Error fetching subreddits:", error);
    }
  };

  const handleIngestion = async (subredditName: string, timeRange: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/subreddits/ingest`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          subreddit: { name: subredditName, sortBy: timeRange },
        }),
      });
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const data = await response.json();
      toast(`Ingesting data from r/${data.Name}`);
    } catch (error) {
      console.error("Error ingesting subreddit:", error);
      toast.error("Failed to ingest subreddit");
    }
  };

  const handleIngestionAll = async (timeRange: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/subreddits/ingest-all`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sortBy: timeRange,
        }),
      });
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      toast(`Ingesting all subreddits for ${timeRange}`);
    } catch (error) {
      console.error("Error ingesting all subreddits:", error);
      toast.error("Failed to ingest all subreddits");
    }
  };

  const handleAddSubreddit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newSubreddit) {
      try {
        console.log(
          `Adding new subreddit: ${newSubreddit} with time range: ${timeRange}`
        );
        setNewSubreddit("");
        const response = await fetch(`${API_BASE_URL}/subreddits/ingest`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            subreddit: { name: newSubreddit, sortBy: timeRange },
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
        toast(`Subscribed to r/${newSubreddit}`);
      } catch (error) {
        console.error("Error adding subreddit:", error);
        toast.error("Failed to add subreddit");
      }
    }
  };

  return (
    <div className="flex flex-col lg:flex-row w-full">
      <div className="w-full lg:w-1/2 p-4">
        <h1 className="text-2xl font-bold mb-4">Subscribed Subreddits</h1>

        {subreddits.length > 0 && (
          <div>
            <div className="mb-8 flex flex-row gap-4">
              <Button variant="secondary" onClick={() => handleIngestionAll("month")}>
                Ingest All Month
              </Button>
              <Button variant="secondary" onClick={() => handleIngestionAll("day")}>
                Ingest All Day
              </Button>
            </div>
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
                  <TableRow 
                    key={subreddit.name}
                    className={selectedSubreddit === subreddit.name ? "bg-muted" : ""}
                    onClick={() => setSelectedSubreddit(subreddit.name)}
                    style={{ cursor: "pointer" }}
                  >
                    <TableCell>{subreddit.name}</TableCell>
                    <TableCell>{subreddit.numberOfSubscribers}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button
                          variant="outline"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleIngestion(subreddit.name, "day");
                          }}
                        >
                          Ingest Day
                        </Button>
                        <Button
                          variant="secondary"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleIngestion(subreddit.name, "month");
                          }}
                        >
                          Ingest Month
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
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
      <SubredditPosts selectedSubreddit={selectedSubreddit} />
    </div>
  );
};

export default SubredditDashboard;
