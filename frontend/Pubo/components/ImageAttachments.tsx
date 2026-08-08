import { Ionicons } from "@expo/vector-icons";
import * as ImagePicker from 'expo-image-picker'
import React,{useState} from 'react';
import { ActivityIndicator,Image,StyleSheet,TouchableOpacity,View } from "react-native";
import { theme } from '@/constants/themes';
import { PostImage } from "@/types/Post";
interface Props{
    images:PostImage[];
    onAdd:(image:PostImage)=>void;
    onRemove:(id:string)=>void;
     // Plug your existing upload endpoint in here — it takes a local file URI
  // and should return the uploaded image's id + public URL.
  onUpload:(asset:{uri:string; fileName?:string | null; mimeType?:string | null; file?: File})=>Promise<{id:string; url:string}>

}
export function ImageAttachments({images,onAdd,onRemove,onUpload}:Props){
    const[isUploading,setIsUploading]=useState(false);
    async function handlePick(){
        const permission=await ImagePicker.requestMediaLibraryPermissionsAsync();
        if(!permission.granted) return;
        const result =await ImagePicker.launchImageLibraryAsync({
            mediaTypes:ImagePicker.MediaTypeOptions.Images,
            allowsMultipleSelection:true,
            quality:0.8,
        });
        if(result.canceled)return;
        setIsUploading(true);
        try{
            for(const asset of result.assets){
                const uploaded=await onUpload(asset);
                onAdd({id:uploaded.id,url:uploaded.url});
            }
        }
        finally{
            setIsUploading(false)
        }
    }
    return (
        <View style={styles.row}>
            {images.map((image)=>(
                <View key={image.id} style={styles.thumbWrap}>
                    <Image source={{uri:image.url}} style={styles.thumb}/>
                    <TouchableOpacity style={styles.removeBtn} onPress={()=>onRemove(image.id)}>
                        <Ionicons name="close" size={14} color='#fff'/>
                    </TouchableOpacity>
                </View>
            ))}
            <TouchableOpacity style={styles.addBtn} onPress={handlePick} disabled={isUploading}>
                {isUploading ?(
                    <ActivityIndicator size='small'/>
                ):(
                    <Ionicons name='add' size={22} color={theme.colors.inkMuted}/>
                )}
            </TouchableOpacity>
        </View>
    )

}
 
const styles = StyleSheet.create({
  row: { flexDirection: 'row', flexWrap: 'wrap', gap: 10, marginTop: 4 },
  thumbWrap: { position: 'relative' },
  thumb: { width: 64, height: 64, borderRadius: theme.radius.sm },
  removeBtn: {
    position: 'absolute',
    top: -6,
    right: -6,
    backgroundColor: theme.colors.ink,
    borderRadius: 10,
    padding: 3,
  },
  addBtn: {
    width: 64,
    height: 64,
    borderRadius: theme.radius.sm,
    borderWidth: 1,
    borderColor: theme.colors.border,
    borderStyle: 'dashed',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: theme.colors.surface,
  },
});
