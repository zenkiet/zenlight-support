export namespace app {
	
	export class UpdateInfo {
	    available: boolean;
	    currentVersion: string;
	    build?: string;
	    latestVersion: string;
	    releaseNotes: string;
	    downloadUrl: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.currentVersion = source["currentVersion"];
	        this.build = source["build"];
	        this.latestVersion = source["latestVersion"];
	        this.releaseNotes = source["releaseNotes"];
	        this.downloadUrl = source["downloadUrl"];
	        this.error = source["error"];
	    }
	}

}

export namespace domain {
	
	export class InstallFileDTO {
	    name: string;
	    data: number[];
	
	    static createFrom(source: any = {}) {
	        return new InstallFileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.data = source["data"];
	    }
	}
	export class ResourceConfig {
	    id: string;
	    name: string;
	    description: string;
	    type: string;
	    path: string;
	    installable: boolean;
	    serviceName?: string;
	
	    static createFrom(source: any = {}) {
	        return new ResourceConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.installable = source["installable"];
	        this.serviceName = source["serviceName"];
	    }
	}
	export class ResourceMetrics {
	    pid: number;
	    createTime: number;
	    cpu: number;
	    mem: number;
	    totalSize?: number;
	    lastModified?: number;
	
	    static createFrom(source: any = {}) {
	        return new ResourceMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.createTime = source["createTime"];
	        this.cpu = source["cpu"];
	        this.mem = source["mem"];
	        this.totalSize = source["totalSize"];
	        this.lastModified = source["lastModified"];
	    }
	}

}

