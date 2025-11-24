import { useState, useEffect } from 'react';
import { SendPrompt } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import FlightRecorder from './components/FlightRecorder';
import './App.css';

interface LogEntry {
    level: string;
    message: string;
    timestamp: number;
}

function App() {
    const [prompt, setPrompt] = useState('');
    const [messages, setMessages] = useState<Array<{ role: string; content: string }>>([]);
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const [isProcessing, setIsProcessing] = useState(false);

    useEffect(() => {
        // Listen for log events from the backend
        EventsOn('kortex:log', (data: { level: string; message: string }) => {
            const logEntry: LogEntry = {
                level: data.level,
                message: data.message,
                timestamp: Date.now(),
            };
            setLogs((prev) => [...prev, logEntry]);

            // If it's a completion or error, stop processing
            if (data.level === 'COMPLETE' || data.level === 'ERROR') {
                setIsProcessing(false);
            }
        });
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!prompt.trim() || isProcessing) return;

        // Add user message to chat
        setMessages((prev) => [...prev, { role: 'user', content: prompt }]);
        setIsProcessing(true);

        try {
            const response = await SendPrompt(prompt);
            setMessages((prev) => [...prev, { role: 'assistant', content: response }]);
            setPrompt('');
        } catch (error) {
            console.error('Failed to send prompt:', error);
            setIsProcessing(false);
        }
    };

    return (
        <div className="app-container">
            {/* Left Panel: Chat Interface */}
            <div className="chat-panel">
                <div className="chat-header">
                    <h1>
                        <span className="logo-icon">⚡</span>
                        KORTEX
                    </h1>
                    <p className="subtitle">Autonomous Interface Layer</p>
                </div>

                <div className="messages-container">
                    {messages.length === 0 ? (
                        <div className="welcome-message">
                            <h2>Welcome to Kortex</h2>
                            <p>Your AI-powered web automation agent.</p>
                            <div className="example-prompts">
                                <p>Try asking:</p>
                                <ul>
                                    <li>"Navigate to google.com"</li>
                                    <li>"Search for AI news"</li>
                                    <li>"Click the first result"</li>
                                </ul>
                            </div>
                        </div>
                    ) : (
                        messages.map((msg, idx) => (
                            <div key={idx} className={`message ${msg.role}`}>
                                <div className="message-content">{msg.content}</div>
                            </div>
                        ))
                    )}
                </div>

                <form onSubmit={handleSubmit} className="input-form">
                    <input
                        type="text"
                        value={prompt}
                        onChange={(e) => setPrompt(e.target.value)}
                        placeholder="Enter your command..."
                        className="prompt-input"
                        disabled={isProcessing}
                    />
                    <button type="submit" className="send-button" disabled={isProcessing || !prompt.trim()}>
                        {isProcessing ? (
                            <span className="loading-spinner">⟳</span>
                        ) : (
                            <span>→</span>
                        )}
                    </button>
                </form>
            </div>

            {/* Right Panel: Mission Control */}
            <div className="mission-control-panel">
                <FlightRecorder logs={logs} />
            </div>
        </div>
    );
}

export default App;
