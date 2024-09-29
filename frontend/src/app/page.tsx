"use client"

import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"

export default function Home() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"

  const [userInput, setUserInput] = useState('')
  const [timeFrame, setTimeFrame] = useState('day')

  const handleSubscription = async () => {
    console.log('Ingesting:', { userInput, timeFrame })
    const response = await fetch(`${apiUrl}/subreddits/ingest`, {
      method: "POST",
      body: JSON.stringify({ subreddits: [{name: userInput, sort_by: timeFrame }]}),
    });
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    const _ = await response.json();
    toast(`Subscribed to r/${userInput}`, {})

    // Here you would typically send this data to your backend
  }
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <Card className="w-full max-w-md mx-auto">
          <CardHeader>
            <CardTitle>Subreddits</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="userInput" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Subreddit name
              </label>
              <Input
                id="userInput"
                value={userInput}
                onChange={(e) => setUserInput(e.target.value)}
                placeholder="travel..."
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="timeFrame" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Top Posts From
              </label>
              <Select value={timeFrame} onValueChange={setTimeFrame}>
                <SelectTrigger id="timeFrame">
                  <SelectValue placeholder="Select" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="day">Day</SelectItem>
                  <SelectItem value="month">Month</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
          <CardFooter>
            <Button onClick={handleSubscription} className="w-full">Subscribe</Button>
          </CardFooter>
        </Card>
      </main>
    </div>
  );
}
