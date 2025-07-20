import { useEffect } from "react";
import ChatPage from "@/pages/ChatPage";
import { OpenFeatureProvider, OpenFeature } from "@openfeature/react-sdk";
import DevCycleReactProvider from "@devcycle/openfeature-react-provider";

function App() {
  useEffect(() => {
    const setupOpenFeature = async () => {
      const devcycleClientSdkKey = import.meta.env.VITE_DEVCYCLE_CLIENT_SDK_KEY;
      const devcycleUserId = import.meta.env.VITE_DEVCYCLE_USER_ID;
      if (!devcycleClientSdkKey || !devcycleUserId) {
        throw new Error(
          "VITE_DEVCYCLE_CLIENT_SDK_KEY or VITE_DEVCYCLE_USER_ID environment variable is not set"
        );
      }
      await OpenFeature.setContext({
        user_id: devcycleUserId,
      });
      await OpenFeature.setProviderAndWait(
        new DevCycleReactProvider(devcycleClientSdkKey)
      );
    };
    setupOpenFeature();
  }, []);

  return (
    <OpenFeatureProvider>
      <ChatPage />
    </OpenFeatureProvider>
  );
}

export default App;
