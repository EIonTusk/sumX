import { useState } from "react";
import { CopyButton } from "@/components/ui/shadcn-io/copy-button";
import { Card, CardAction, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

type Props = {
  data: ReqObj;
};

type Summary = {
    heading: string
    text: string
}


type ReqObj = {
    params: {
        username: string
        from: string
        to: string
        limit: number
    }
    summary: Summary[]
    tweets: string[]
}

export default function SummaryCard({ data }: Props) {
    const [showRaw, setShowRaw] = useState(false);

    return (
        <Card className="w-full bg-card max-w-2xl min-w-2xl shadow-2xl rounded-2xl backdrop-blur ml-24 mb-24">
            <CardHeader>
                <CardTitle>
                    {data.params.username}
                </CardTitle>
                <CardDescription>
                    {data.params.limit != -1 ? `last ${data.params.limit} posts ` : ""}
                    {data.params.from != "" ? `from ${data.params.from} ` : ""}
                    {data.params.to != "" ? `to ${data.params.to} ` : ""}
                </CardDescription>
                <CardAction>
                    <CopyButton content={data.summary.map(item  => `${item.heading}\n${item.text}`).join("\n\n")} size="md" variant="secondary" />
                </CardAction>
            </CardHeader>
            <CardContent className="">
                <div className="border border-input p-4 rounded-md text-sm space-y-2" id="summary">
                      {data.summary.map((item, index) => {
                          return (
                        <div key={index}>
                            <h3 className="font-bold">{item.heading}</h3>
                            <p className="">{item.text}</p>
                        </div>);
                      })}
                </div>
                <div className="flex justify-center my-4">
                    <button
                        onClick={() => setShowRaw(!showRaw)}
                        className="text-sm hover:underline text-muted-foreground cursor-pointer"
                    >
                        {showRaw ? "Hide raw tweets" : "Show raw tweets"}
                    </button>
                </div>
                {showRaw && (
                    <div className="border border-input p-4 rounded-md text-sm space-y-2" id="summary">
                          {data.tweets.map((item, index) => {
                              return (
                            <div key={index}>
                                <p className="">{item}</p>
                            </div>);
                          })}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
