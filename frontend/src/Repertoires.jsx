import React, { useEffect, useState } from "react";
import * as RepMgr from "../wailsjs/go/backend/RepertoireManager";

export default function Repertoires({ onSelect, onResetFen, setRefreshRepertoires, setFen, onPracticeModeChange }) {
    const [items, setItems] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [selectedId, setSelectedId] = useState(null);
    const [practiceFEN, setPracticeFEN] = useState(null);
    const [practiceMove, setPracticeMove] = useState("");

    const [name, setName] = useState("");
    const [color, setColor] = useState("white");
    const [elo, setElo] = useState("1200");
    const [coverage, setCoverage] = useState(50);

    const eloOptions = ["0", "1200", "1400", "1600", "1800", "2000+"];
    const coverageOptions = [50, 100, 150, 200];

    async function refresh() {
        setLoading(true);
        setError(null);
        try {
            const list = await RepMgr.List();
            const mapped = await Promise.all(
                list.map(async (r) => {
                    const dueCount = await RepMgr.CountDueNodes(r.id); // Fetch due nodes count
                    return {
                        id: r.id,
                        name: r.name,
                        color: r.color === "white" ? "white" : "black",
                        elo: String(r.elo),
                        coverage: r.coverage,
                        dueCount, // Add due nodes count to the repertoire object
                    };
                })
            );
            setItems(mapped);
        } catch (e) {
            setError(String(e));
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        refresh();
        setRefreshRepertoires(() => refresh); // Provide refresh function to parent
    }, []);

    async function create() {
        try {
            await RepMgr.Create(name.trim(), color, parseInt(elo));
            setName("");
            setColor("white");
            setElo("1200");
            setCoverage(50);
            await refresh();
        } catch (e) {
            setError(String(e));
        }
    }

    async function update(r) {
        try {
            await RepMgr.Update({ ...r, elo: parseInt(r.elo) });
            await refresh();
        } catch (e) {
            setError(String(e));
        }
    }

    async function remove(id) {
        if (!window.confirm("Delete repertoire?")) return;
        try {
            await RepMgr.Delete(id);
            if (id === selectedId) {
                setSelectedId(null); // Unselect if the deleted repertoire was selected
                if (onSelect) onSelect(null); // Notify parent of unselection
                if (onResetFen) onResetFen(); // Reset chessboard to starting position
            }
            await refresh();
        } catch (e) {
            setError(String(e));
        }
    }

    async function select(id) {
        try {
            await RepMgr.SelectRepertoire(id); // ✅ backend state
            setSelectedId(id); // ✅ frontend state
            if (onSelect) onSelect(id); // ✅ notify parent
            if (onResetFen) onResetFen(); // Reset chessboard to starting position
        } catch (e) {
            setError(String(e));
        }
    }

    async function startPractice(repID) {
        try {
            const dueFENs = await RepMgr.GetDueFENs();
            if (!dueFENs || dueFENs.length === 0) {
                alert("No due positions available for practice.");
                return;
            }
            const fen = dueFENs[0]; // Use the first due position
            setPracticeFEN(fen);
            setFen(fen);
            RepMgr.SetCurrentFEN(fen);
            if (onPracticeModeChange) onPracticeModeChange(true); // Notify parent to hide MoveList
        } catch (e) {
            setError(`Failed to fetch due positions: ${e}`);
        }
    }

    async function submitPracticeMove() {
        try {
            await RepMgr.TestCurrentPositionWithDueDate(practiceMove);
            alert("Move is correct!");
        } catch (e) {
            alert(`Incorrect move: ${e.message}`);
        } finally {
            setPracticeFEN(null); // Close practice window
            setPracticeMove(""); // Reset move input
            const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"; // Define starting position
            setFen(startingFEN); // Reset chessboard to starting position
            await RepMgr.SetCurrentFEN(startingFEN); // Update backend to starting position
            if (onPracticeModeChange) onPracticeModeChange(false); // Notify parent to show MoveList
            await refresh(); // Refresh the repertoire list
        }
    }

    return (
        <div style={{ color: "green" }}>
            {loading && <p>Loading...</p>}
            {error && <p style={{ color: "red" }}>{error}</p>}

            <div>
                <h3>Create new</h3>
                <input
                    placeholder="Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <select value={color} onChange={(e) => setColor(e.target.value)}>
                    <option value="white">white</option>
                    <option value="black">black</option>
                </select>
                <select value={elo} onChange={(e) => setElo(e.target.value)}>
                    {eloOptions.map((opt) => (
                        <option key={opt} value={opt}>{opt}</option>
                    ))}
                </select>
                <select value={coverage} onChange={(e) => setCoverage(Number(e.target.value))}>
                    {coverageOptions.map((opt) => (
                        <option key={opt} value={opt}>{opt}</option>
                    ))}
                </select>
                <button onClick={create} disabled={!name.trim()}>Create</button>
            </div>

            {items.map((r) => (
                <div
                    key={r.id}
                    style={{
                        border: "1px solid #ddd",
                        margin: "0.5rem",
                        padding: "0.5rem",
                        background: r.id === selectedId ? "#eef" : "transparent"
                    }}
                >
                    <input
                        value={r.name}
                        onChange={(e) =>
                            setItems((prev) =>
                                prev.map((x) => (x.id === r.id ? { ...x, name: e.target.value } : x))
                            )
                        }
                    />
                    <input value={r.color} disabled />
                    <select
                        value={r.elo}
                        onChange={(e) =>
                            setItems((prev) =>
                                prev.map((x) => (x.id === r.id ? { ...x, elo: e.target.value } : x))
                            )
                        }
                    >
                        {eloOptions.map((opt) => (
                            <option key={opt} value={opt}>{opt}</option>
                        ))}
                    </select>
                    <select
                        value={r.coverage}
                        onChange={(e) =>
                            setItems((prev) =>
                                prev.map((x) => (x.id === r.id ? { ...x, coverage: Number(e.target.value) } : x))
                            )
                        }
                    >
                        {coverageOptions.map((opt) => (
                            <option key={opt} value={opt}>{opt}</option>
                        ))}
                    </select>
                    <p>{r.dueCount} Due Nodes</p> {/* Display due nodes count */}
                    <button onClick={() => update(r)}>Save</button>
                    <button onClick={() => select(r.id)}>Select</button>
                    <button onClick={() => remove(r.id)}>Delete</button>
                    <button onClick={() => startPractice(r.id)} disabled={selectedId !== r.id}>Practice</button> {/* Practice button */}
                </div>
            ))}

            {practiceFEN && (
                <div style={{ marginTop: "1rem", padding: "1rem", border: "1px solid #ccc" }}>
                    <h4>Practice Mode</h4>
                    <p>Current Position: {practiceFEN}</p>
                    <input
                        placeholder="Enter your move (SAN)"
                        value={practiceMove}
                        onChange={(e) => setPracticeMove(e.target.value)}
                    />
                    <button onClick={submitPracticeMove}>Submit Move</button>
                </div>
            )}
        </div>
    );
}
