import React, { useState } from "react";
import {
  ActivityIndicator,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";
import { CharacterMeter } from "./characterMeter";
import { ImageAttachments } from "../components/ImageAttachments";
import { PlatformSelector } from "../components/PlatformSelector";
import { ScheduleModal } from "./scheduleModal";
import { exceedsLimit } from "@/constants/Platforms";
import { theme } from "../constants/themes";
import { usePostComposer } from "../hooks/usePostComposer";
import { Ionicons } from "@expo/vector-icons";
import { useAuth } from "@clerk/clerk-expo";
import { postsApi } from "@/api/Postsapi";

interface Props {
  // Optional: pass this when the screen is shown in a Modal/sheet so there's
  // a way to dismiss it. Safe to omit if this is a normal routed screen.
  onClose?: () => void;
}

export function ComposePostScreen({ onClose }: Props) {
  const composer = usePostComposer();
  const { getToken } = useAuth();
  const [isScheduleOpen, setIsScheduleOpen] = useState(false);

  const overLimit = exceedsLimit(composer.content, composer.platforms);
  const disableSubmit = !composer.canSubmit || overLimit || composer.isSaving;

  async function handleUpload(asset: { uri: string; fileName?: string | null; mimeType?: string | null; file?: File }) {
    const token = await getToken();
    if (!token) throw new Error("Please sign in before uploading an image");
    const uploaded = await postsApi.uploadImage(token, asset.uri, asset.fileName ?? "image.jpg", asset.mimeType ?? "image/jpeg", asset.file);
    return { id: uploaded.id, url: uploaded.storageURL };
  }

  return (
    <ScrollView style={styles.screen} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <Text style={styles.heading}>New post</Text>
        {onClose && (
          <TouchableOpacity onPress={onClose} hitSlop={10}>
            <Ionicons name="close" size={24} color={theme.colors.inkMuted} />
          </TouchableOpacity>
        )}
      </View>

      <TextInput
        style={styles.input}
        multiline
        placeholder="What did you build today?"
        placeholderTextColor={theme.colors.inkMuted}
        value={composer.content}
        onChangeText={composer.setContent}
      />

      <CharacterMeter
        content={composer.content}
        platforms={composer.platforms}
      />

      <Text style={styles.sectionLabel}>Post to</Text>
      <PlatformSelector
        selected={composer.platforms}
        onToggle={composer.togglePlatform}
      />

      <Text style={styles.sectionLabel}>Images</Text>
      <ImageAttachments
        images={composer.images}
        onAdd={composer.addImage}
        onRemove={composer.removeImage}
        onUpload={handleUpload}
      />

      {composer.error && <Text style={styles.errorText}>{composer.error}</Text>}

      <View style={styles.footer}>
        <TouchableOpacity
          style={[styles.actionBtn, styles.draftBtn]}
          onPress={async () => { const token = await getToken(); if (!token) return; await composer.saveDraft(token); }}
          disabled={composer.isSaving}
        >
          <Text style={styles.draftText}>Save draft</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[
            styles.actionBtn,
            styles.scheduleBtn,
            disableSubmit && styles.disabled,
          ]}
          onPress={() => setIsScheduleOpen(true)}
          disabled={disableSubmit}
        >
          <Text style={styles.scheduleText}>Schedule</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[
            styles.actionBtn,
            styles.postBtn,
            disableSubmit && styles.disabled,
          ]}
          onPress={async () => { const token = await getToken(); if (!token) throw new Error("Please sign in before posting"); await composer.publishNow(token); }}
          disabled={disableSubmit}
        >
          {composer.isSaving ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.postText}>Post now</Text>
          )}
        </TouchableOpacity>
      </View>

      <ScheduleModal
        visible={isScheduleOpen}
        onClose={() => setIsScheduleOpen(false)}
        onConfirm={async (iso) => {
          await composer.schedule(iso);
          setIsScheduleOpen(false);
        }}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: theme.colors.background },
  content: { padding: 20, paddingBottom: 40 },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 16,
  },
  heading: {
    fontSize: theme.font.title,
    fontWeight: "700",
    color: theme.colors.ink,
  },
  input: {
    minHeight: 120,
    fontSize: theme.font.body,
    color: theme.colors.ink,
    backgroundColor: theme.colors.surface,
    borderRadius: theme.radius.md,
    borderWidth: 1,
    borderColor: theme.colors.border,
    padding: 14,
    textAlignVertical: "top",
  },
  sectionLabel: {
    fontSize: theme.font.label,
    color: theme.colors.inkMuted,
    fontWeight: "700",
    marginTop: 20,
    marginBottom: 8,
    textTransform: "uppercase",
  },
  errorText: { color: theme.colors.danger, marginTop: 12 },
  footer: { flexDirection: "row", gap: 10, marginTop: 28 },
  actionBtn: {
    flex: 1,
    paddingVertical: 14,
    borderRadius: theme.radius.md,
    alignItems: "center",
    justifyContent: "center",
  },
  draftBtn: {
    backgroundColor: theme.colors.surface,
    borderWidth: 1,
    borderColor: theme.colors.border,
  },
  draftText: { color: theme.colors.ink, fontWeight: "700" },
  scheduleBtn: { backgroundColor: theme.colors.accentMuted },
  scheduleText: { color: theme.colors.accent, fontWeight: "700" },
  postBtn: { backgroundColor: theme.colors.accent },
  postText: { color: "#fff", fontWeight: "700" },
  disabled: { opacity: 0.5 },
});
