import { useState, useRef, useEffect } from "react";
import { ThemeProvider } from "@/components/theme-provider"
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, } from "@/components/ui/card";
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { toast, Toaster } from "sonner";
import { Loader2 } from "lucide-react";
import { parseISO, formatDuration, intervalToDuration } from 'date-fns';
import SummaryCard, { type ReqObj} from "@/components/summary-card";


function App() {
    const [username, setUsername] = useState("");
    const [limit, setLimit] = useState("");
    const [retweets, setRetweets] = useState(false)
    const [fromDate, setFromDate] = useState("");
    const [toDate, setToDate] = useState("");
    const [showAdvanced, setShowAdvanced] = useState(false);
    const [showLast, setShowLast] = useState(false);
    const [loading, setLoading] = useState(false);
    const [timeLeft, setTimeLeft] = useState("");
    const [summaries, setSummaries] = useState<ReqObj[]>([]);
    const [delta, setDelta] = useState('');

    const ref = useRef<HTMLDivElement>(null);

    const today = new Date().toISOString().split('T')[0];

    const fetchData = async () => {
        if (!username) return;
        setLoading(true);

        const query = new URLSearchParams({});
        if (limit) query.append("limit", limit);
        if (retweets) query.append("retweets", "1")
        if (fromDate) query.append("from", fromDate+"T00:00:00Z");
        if (toDate) query.append("to", toDate+"T23:59:59Z");

        try {
            const res = await fetch(`${import.meta.env.VITE_API_URL}/api/summarize/${username}${query.size !== 0 ? "?" : ""}${query.toString()}`);

            const data = await res.json();
            console.log(data);

            if (!res.ok) {
                if (res.status == 429) {
                    setTimeLeft(data.next_reset)
                    toast.error("Request limit reached", {
                        description: data.error,
                    });
                    return
                }
                toast.error("Error occured", {
                    description: data.error,
                });
                return;
            }

            if (data.data.params.username.at(0) !== '@') {
                data.data.params.username = '@' + data.data.params.username;
            }
            summaries.unshift(data.data);
            setShowLast(true);
        } catch (err) {
            console.error("Failed to fetch data", err);
            toast.error("Error occured", {
                description: "Failed to fetch data",
            });
        } finally {
            setLoading(false);
        }
    };

    function userInput(e: React.ChangeEvent<HTMLInputElement>) {
        const username: string = e.target.value;
        if (username === "@") {
            setUsername("");
            return;
        }
        if (username.at(0) === "@" || username.length === 0) {
            setUsername(username);
            return;
        }
        setUsername("@"+username);
    }

    useEffect(() => {
        const el: any = ref.current;
        if (!el) return;

        // Manually trigger transition by changing styles
        if (showAdvanced) {
            // Set height to current height first
            el.style.height = el.scrollHeight + 'px';
        } else {
            // Collapse height to 0
            el.style.height = '376px';
        }
    }, [showAdvanced]);

    useEffect(() => {
        const fetchSummaries = async () => {
            try {
                const res = await fetch(`${import.meta.env.VITE_API_URL}/api/summaries`);
                const data = await res.json();
                console.log(data)
                
                if (!res.ok) {
                    toast.error("Error occured", {
                        description: data.error,
                    });
                    return;
                }

                if (!data.data) {
                    return
                }

                const cleanedData = data.data.map((e: ReqObj) => {
                    return {
                        ...e,
                        params: {
                            ...e.params,
                            username: e.params.username.startsWith("@")
                                ? e.params.username
                                : "@" + e.params.username,
                        },
                    };
                });

                setSummaries(cleanedData);
            } catch (err) {
                console.error("Failed to fetch data", err);
                toast.error("Error occured", {
                    description: "Failed to fetch data",
                });
            } finally {
                setLoading(false);
            }
        }
        fetchSummaries()
    }, []);

      useEffect(() => {
        const target = parseISO(timeLeft);

        const updateDelta = () => {
          const now = new Date();
          const duration = intervalToDuration({
            start: now < target ? now : target,
            end: now < target ? target : now,
          });

          const formatted = formatDuration(duration, { delimiter: ', ' }).replace(/\bminutes?\b/g, 'min')
              .replace(/\bhours?\b/g, 'h')
              .replace(/\bdays?\b/g, 'd')
              .replace(/\bseconds?\b/g, 's');
          if (now > target) {
              setDelta("");
              return;
          }
          setDelta(`${formatted}`);
        };

        updateDelta(); // Initial run
        const interval = setInterval(updateDelta, 1000);

        return () => clearInterval(interval);
      }, [timeLeft]);


    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <Toaster />
            <div className="max-h-screen min-h-screen relative flex items-center justify-center p-4 overflow-hidden min-w-screen">
                <div className="absolute right-0 top-0 m-8 z-10">
                    <Button onClick={() => {setShowLast(!showLast)}} className="backdrop-blur-md bg-primary/60 cursor-pointer">
                        {showLast ? "hide summaries" : "show summaries"}
                    </Button>
                </div>
                <div className="flex flex-col max-h-screen items-center justify-center p-4 overflow-hidden">
                    <Card
                    className="overflow-hidden transition-[height] duration-500 ease-in-out w-full bg-card max-w-2xl min-w-2xl shadow-2xl rounded-2xl backdrop-blur" id="card"
                    ref={ref}
                    >
                    <CardContent className="p-8 space-y-6">
                        <h1 className="text-3xl font-extrabold text-center">sumX</h1>
                        <p className="text-center text-muted-foreground text-sm">Enter a X username to generate a summary of recent tweets</p>
                        <div className="flex md:flex-row md:space-x-2 items-center justify-center space-y-2 md:space-y-0">
                            <Input
                                placeholder="@username"
                                value={username}
                                onChange={userInput}
                                className="flex-1 max-w-sm text-center rounded-xl shadow"
                                />
                        </div>

                        <div className="flex justify-center">
                            <button
                                onClick={() => setShowAdvanced(!showAdvanced)}
                                className="text-sm hover:underline text-muted-foreground cursor-pointer"
                            >
                                {showAdvanced ? "Hide advanced options" : "Show advanced options"}
                            </button>
                        </div>

                        {(showAdvanced) && (
                            <>
                                <div className="flex flex-col md:flex-row md:space-x-2 items-center justify-center space-y-2 md:space-y-0">
                                    <div className="grid w-full max-w-sm items-center gap-3">
                                        <Label htmlFor="from">From</Label>
                                        <Input
                                            placeholder={today}
                                            value={fromDate}
                                            onChange={(e) => setFromDate(e.target.value)}
                                            type="date"
                                            className="flex-1 max-w-sm text-center rounded-xl shadow"
                                        />
                                    </div>
                                    <div className="grid w-full max-w-sm items-center gap-3">
                                        <Label htmlFor="to">To</Label>
                                        <Input
                                            placeholder={today}
                                            value={toDate}
                                            onChange={(e) => setToDate(e.target.value)}
                                            type="date"
                                            className="flex-1 max-w-sm text-center rounded-xl shadow"
                                        />
                                    </div>
                                    <div className="grid w-full max-w-sm items-center gap-3">
                                        <Label htmlFor="limit">Limit</Label>
                                        <Input
                                            placeholder="10"
                                            value={limit}
                                            onChange={(e) => setLimit(e.target.value)}
                                            type="number"
                                            min="5"
                                            max="100"
                                            className="flex-1 max-w-sm text-center rounded-xl shadow"
                                        />
                                    </div>
                                </div>
                                <div className="flex justify-center items-center align-middle justify-items-center text-center gap-3 w-full">
                                    <Checkbox checked={retweets} onCheckedChange={() => setRetweets(!retweets)}/>
                                    <p className="text-sm text-muted-foreground">include retweets</p>
                                </div>
                            </>
                        )}

                        <div className="flex justify-center">
                            <Button variant={delta === "" ? "default" : "ghost"} onClick={fetchData} disabled={loading || delta !== "" || username === ""} className="rounded-xl px-6 cursor-pointer backdrop-blur-md bg-primary/60">
                                {loading ? <Loader2 className="animate-spin" /> : delta !== "" ? <div className="text-foreground">Next request possible in {delta}</div> : "Fetch"}
                            </Button>
                        </div>

                    </CardContent>
                </Card>
            </div>
            <div className="flex flex-col w-max min-h-screen max-h-screen overflow-y-auto pt-24 [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
            {(showLast) && (
                <div className="min-w-2xl">
                </div>
            )}
            {(showLast) && (
                summaries.map((v, i) => {
                    return (
                        <SummaryCard key={i} data={v} />
                    )
                })
            )}
            </div>
        </div>
    </ThemeProvider>
    );
}

export default App;
