
export interface TestMessage {
    TypTestMessage?: unknown;
    Id: number;
    Name: string;
    Ok: boolean;
}
export function isTestMessage(obj: any): obj is TestMessage {
    return typeof obj === 'object' && obj !== null && 'TypTestMessage' in obj;
}
export interface Test1Message {
    TypTest1Message?: unknown;
    Id1: number;
    Name1: string;
    Ok1: boolean;
}
export function isTest1Message(obj: any): obj is Test1Message {
    return typeof obj === 'object' && obj !== null && 'TypTest1Message' in obj;
}
