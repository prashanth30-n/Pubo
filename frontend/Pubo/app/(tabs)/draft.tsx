import React, { useEffect, useRef, useState } from 'react';
import { ActivityIndicator, FlatList, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuth } from '@clerk/clerk-expo';
import useTheme from '@/hooks/useTheme';
import { postsApi, SavedPost } from '@/api/Postsapi';

export default function Drafts() {
  const { colors } = useTheme(); const { getToken } = useAuth(); const getTokenRef=useRef(getToken); const [posts,setPosts]=useState<SavedPost[]>([]); const [loading,setLoading]=useState(true); const [error,setError]=useState<string | null>(null);
  useEffect(()=>{getTokenRef.current=getToken},[getToken]);
  useEffect(()=>{(async()=>{try{const token=await getTokenRef.current();if(token)setPosts(await postsApi.list(token,'draft'));}catch(e){setError(e instanceof Error?e.message:'Could not load drafts');}finally{setLoading(false)}})()},[]);
  return <SafeAreaView style={[styles.screen,{backgroundColor:colors.backGround}]}><Text style={[styles.title,{color:colors.text}]}>Drafts</Text>{loading?<ActivityIndicator/>:error?<Text style={styles.error}>{error}</Text>:<FlatList data={posts} keyExtractor={(post)=>post.id} ListEmptyComponent={<Text style={[styles.empty,{color:colors.text}]}>No drafts yet.</Text>} renderItem={({item})=><View style={styles.card}><Text style={styles.content}>{item.content || 'Image draft'}</Text><Text style={styles.meta}>{item.mediaIds.length} image{item.mediaIds.length===1?'':'s'} · {new Date(item.createdAt).toLocaleDateString()}</Text></View>}/>}</SafeAreaView>
}
const styles=StyleSheet.create({screen:{flex:1,padding:20},title:{fontSize:24,fontWeight:'700',marginBottom:16},card:{backgroundColor:'#fff',borderRadius:12,padding:14,marginBottom:10},content:{fontSize:16,color:'#172033'},meta:{marginTop:8,color:'#687385'},empty:{textAlign:'center',marginTop:32},error:{color:'#c43b3b'}});
