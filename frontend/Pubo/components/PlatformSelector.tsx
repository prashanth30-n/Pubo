import React from "react";
import { StyleSheet, View, Text, TouchableOpacity } from "react-native";
import { Ionicons, MaterialCommunityIcons } from "@expo/vector-icons";
import { PLATFORM_LIST } from "@/constants/Platforms";
import { theme } from "@/constants/themes";
import { PlatformId } from "@/types/Post";

interface Props{
    selected:PlatformId[];
    onToggle:(id:PlatformId)=>void;
}
export function PlatformSelector({selected,onToggle}:Props){
    return(
        <View style={styles.row}>
            {PLATFORM_LIST.map((platform)=>{
                const isActive=selected.includes(platform.id);
                return(
                    <TouchableOpacity key={platform.id} onPress={()=>onToggle(platform.id)} style={[styles.chip,isActive && {backgroundColor:platform.color,borderColor:platform.color}]} accessibilityRole="checkbox" accessibilityState={{ checked: isActive }}>
                        <View style={[styles.iconBadge, { backgroundColor: isActive ? "#fff" : platform.color }]}>
                          {platform.id === "linkedin" ? <Ionicons name="logo-linkedin" size={18} color={isActive ? platform.color : "#fff"} /> : <MaterialCommunityIcons name="cloud-outline" size={19} color={isActive ? platform.color : "#fff"} />}
                        </View>
                        <Text style={[styles.label,isActive && styles.labelActive]}>{platform.label}</Text>
                    </TouchableOpacity>
                )
            })}
        </View>
    )

}

const styles = StyleSheet.create({
  row: { flexDirection: 'row', flexWrap: 'wrap', gap: 8 },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    paddingVertical: 9,
    paddingHorizontal: 12,
    borderRadius: theme.radius.lg,
    borderWidth: 1,
    borderColor: theme.colors.border,
    backgroundColor: theme.colors.surface,
  },
  iconBadge: { width: 24, height: 24, borderRadius: 12, alignItems: 'center', justifyContent: 'center' },
  label: { fontSize: theme.font.label, color: theme.colors.ink, fontWeight: '600' },
  labelActive: { color: '#fff' },
});
