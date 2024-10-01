"use client";

import SubredditDashboard from "@/components/SubredditDashboard";

export default function Home() {
  return (
    <div className="flex items-center justify-items-center p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        {/* <Card className="w-full max-w-md mx-auto">
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
        </Card> */}

        <SubredditDashboard />
      </main>
    </div>
  );
}
