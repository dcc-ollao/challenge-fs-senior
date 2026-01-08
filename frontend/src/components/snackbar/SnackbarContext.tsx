import { createContext, useContext, useState } from "react";

type SnackbarType = "error" | "success";

type SnackbarState = {
  message: string;
  type: SnackbarType;
} | null;

type SnackbarContextValue = {
  showError: (message: string) => void;
  showSuccess: (message: string) => void;
};

const SnackbarContext = createContext<SnackbarContextValue | null>(null);

export function SnackbarProvider({ children }: { children: React.ReactNode }) {
  const [snackbar, setSnackbar] = useState<SnackbarState>(null);

  function showError(message: string) {
    setSnackbar({ message, type: "error" });
    autoHide();
  }

  function showSuccess(message: string) {
    setSnackbar({ message, type: "success" });
    autoHide();
  }

  function autoHide() {
    setTimeout(() => setSnackbar(null), 3000);
  }

  return (
    <SnackbarContext.Provider value={{ showError, showSuccess }}>
      {children}
      {snackbar && <Snackbar {...snackbar} />}
    </SnackbarContext.Provider>
  );
}

export function useSnackbar() {
  const ctx = useContext(SnackbarContext);
  if (!ctx) throw new Error("useSnackbar must be used within SnackbarProvider");
  return ctx;
}

function Snackbar({ message, type }: { message: string; type: SnackbarType }) {
  return (
    <div className="fixed bottom-6 left-1/2 z-50 -translate-x-1/2">
      <div
        className={[
          "rounded-md px-4 py-2 text-sm shadow-lg",
          type === "error"
            ? "bg-red-600 text-white"
            : "bg-emerald-600 text-white",
        ].join(" ")}
      >
        {message}
      </div>
    </div>
  );
}
