import { useEffect, useRef } from 'react';
import './FlightRecorder.css';

interface LogEntry {
    level: string;
    message: string;
    timestamp: number;
}

interface FlightRecorderProps {
    logs: LogEntry[];
}

function FlightRecorder({ logs }: FlightRecorderProps) {
    const terminalRef = useRef<HTMLDivElement>(null);

    // Auto-scroll to bottom when new logs arrive
    useEffect(() => {
        if (terminalRef.current) {
            terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
        }
    }, [logs]);

    const getLogColor = (level: string): string => {
        switch (level.toUpperCase()) {
            case 'NAVIGATE':
                return 'log-navigate';
            case 'CLICK':
                return 'log-click';
            case 'TYPE':
                return 'log-type';
            case 'HIGHLIGHT':
                return 'log-highlight';
            case 'GET_SNAPSHOT':
                return 'log-snapshot';
            case 'PLANNING':
                return 'log-planning';
            case 'INIT':
                return 'log-init';
            case 'USER':
                return 'log-user';
            case 'ERROR':
                return 'log-error';
            case 'COMPLETE':
                return 'log-complete';
            case 'SHUTDOWN':
                return 'log-shutdown';
            default:
                return 'log-default';
        }
    };

    const formatTimestamp = (timestamp: number): string => {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('en-US', {
            hour12: false,
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    };

    return (
        <div className="flight-recorder">
            <div className="terminal-header">
                <span className="terminal-title">⚙ MISSION CONTROL</span>
                <div className="terminal-indicators">
                    <span className="indicator active"></span>
                    <span className="indicator-label">LIVE</span>
                </div>
            </div>

            <div className="terminal-body" ref={terminalRef}>
                {logs.length === 0 ? (
                    <div className="terminal-empty">
                        <p>Awaiting agent activity...</p>
                        <p className="terminal-hint">Logs will appear here in real-time</p>
                    </div>
                ) : (
                    logs.map((log, idx) => (
                        <div key={idx} className={`terminal-line ${getLogColor(log.level)}`}>
                            <span className="log-timestamp">[{formatTimestamp(log.timestamp)}]</span>
                            <span className="log-level">[{log.level}]</span>
                            <span className="log-message">{log.message}</span>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}

export default FlightRecorder;
