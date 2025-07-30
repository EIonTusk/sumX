import React, { useEffect, useState } from 'react';
import { Card, CardContent } from "@/components/ui/card";
import { differenceInSeconds, parseISO, formatDuration, intervalToDuration } from 'date-fns';

type Props = {
  targetTime: string; // RFC3339 time string
};

export default function TimeDeltaClock({ targetTime }: Props) {

  return (
    <div className="text-muted-foreground">Next request possible in {delta}</div>
  );
}
